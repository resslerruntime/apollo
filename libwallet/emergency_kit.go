package libwallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/muun/libwallet/emergencykit"
)

// EKInput input struct to fill the PDF
type EKInput struct {
	FirstEncryptedKey  string
	FirstFingerprint   string
	SecondEncryptedKey string
	SecondFingerprint  string
	SwapperPublicKey   string
}

// EKOutput with the html as string and the verification code
type EKOutput struct {
	HTML             string
	VerificationCode string
	Metadata         string
}

// GenerateEmergencyKitHTML returns the translated html as a string along with the verification
// code and the kit metadata, represented in an opaque string.
// After calling this method, clients should use their Chromium/WebKit implementations to render
// the HTML into a PDF (better done there), and then come back to call `AddEmergencyKitMetadata`
// and produce the final PDF (better done here).
func GenerateEmergencyKitHTML(ekParams *EKInput, language string) (*EKOutput, error) {
	moduleInput := &emergencykit.Input{
		FirstEncryptedKey:  ekParams.FirstEncryptedKey,
		FirstFingerprint:   ekParams.FirstFingerprint,
		SecondEncryptedKey: ekParams.SecondEncryptedKey,
		SecondFingerprint:  ekParams.SecondFingerprint,
	}

	// Create the HTML and the verification code:
	htmlWithCode, err := emergencykit.GenerateHTML(moduleInput, language)
	if err != nil {
		return nil, fmt.Errorf("GenerateEkHtml failed to render: %w", err)
	}

	// Create and serialize the metadata:
	metadata, err := createEmergencyKitMetadata(ekParams)
	if err != nil {
		return nil, fmt.Errorf("GenerateEkHtml failed to create metadata: %w", err)
	}

	metadataBytes, err := json.Marshal(&metadata)
	if err != nil {
		return nil, fmt.Errorf("GenerateEkHtml failed to marshal %s: %w", string(metadataBytes), err)
	}

	output := &EKOutput{
		HTML:             htmlWithCode.HTML,
		VerificationCode: htmlWithCode.VerificationCode,
		Metadata:         string(metadataBytes),
	}

	return output, nil
}

// AddEmergencyKitMetadata produces a copy of the PDF file at `srcFile` with embedded metadata,
// writing it into `dstFile`. The provided metadata must be the same opaque string produced by
// `GenerateEmergencyKitHTML`.
func AddEmergencyKitMetadata(metadataText string, srcFile string, dstFile string) error {
	// Initialize the MetadataWriter:
	metadataWriter := &emergencykit.MetadataWriter{
		SrcFile: srcFile,
		DstFile: dstFile,
	}

	// Deserialize the metadata:
	var metadata emergencykit.Metadata

	err := json.Unmarshal([]byte(metadataText), &metadata)
	if err != nil {
		return fmt.Errorf("AddEkMetadata failed to unmarshal: %w", err)
	}

	err = metadataWriter.WriteMetadata(&metadata)
	if err != nil {
		return fmt.Errorf("AddEkMetadata failed to write metadata: %w", err)
	}

	return nil
}

func createEmergencyKitMetadata(ekParams *EKInput) (*emergencykit.Metadata, error) {
	// NOTE:
	// This method would be more naturally placed in the `emergencykit` module, but given the current
	// project structure (heavily determined by `gomobile` and the need for top-level bindings) and
	// the use of `decodeEncryptedPrivateKey` this isn't possible. Instead, we peek through the layer
	// boundary to craft the object here.

	// Decode both keys, to extract their inner properties:
	firstKey, err := decodeEncryptedPrivateKey(ekParams.FirstEncryptedKey)
	if err != nil {
		return nil, fmt.Errorf("createEkMetadata failed to decode first key: %w", err)
	}

	secondKey, err := decodeEncryptedPrivateKey(ekParams.SecondEncryptedKey)
	if err != nil {
		return nil, fmt.Errorf("createEkMetadata failed to decode second key: %w", err)
	}

	// Obtain the list of checksumed output descriptors:
	descriptors := emergencykit.GetDescriptors(&emergencykit.DescriptorsData{
		FirstFingerprint:  ekParams.FirstFingerprint,
		SecondFingerprint: ekParams.SecondFingerprint,
	})

	// Create the keys for the key array:
	keys := []*emergencykit.MetadataKey{
		createEmergencyKitMetadataKey(firstKey),
		createEmergencyKitMetadataKey(secondKey),
	}

	// Create the public key array:
	publicKeys := []*emergencykit.MetadataPublicKey{
		{Tag: emergencykit.SwapServerTag, PublicKey: ekParams.SwapperPublicKey},
	}

	metadata := &emergencykit.Metadata{
		Version:           3,
		BirthdayBlock:     int(secondKey.Birthday),
		EncryptedKeys:     keys,
		PublicKeys:        publicKeys,
		OutputDescriptors: descriptors,
	}

	return metadata, nil
}

func createEmergencyKitMetadataKey(key *encryptedPrivateKey) *emergencykit.MetadataKey {
	return &emergencykit.MetadataKey{
		DhPubKey:         hex.EncodeToString(key.EphPublicKey),
		EncryptedPrivKey: hex.EncodeToString(key.CipherText),
		Salt:             hex.EncodeToString(key.Salt),
	}
}

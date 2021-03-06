package io.muun.common.api;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonInclude;

import javax.validation.constraints.NotNull;

@JsonInclude(JsonInclude.Include.NON_NULL)
@JsonIgnoreProperties(ignoreUnknown = true)
public class CreateFirstSessionOkJson {

    @NotNull
    public PublicKeyJson cosigningPublicKey;

    @NotNull
    public PublicKeyJson swapServerPublicKey;

    @NotNull
    public UserJson user;

    /**
     * Json constructor.
     */
    public CreateFirstSessionOkJson() {
    }

    /**
     * Manual constructor.
     */
    public CreateFirstSessionOkJson(UserJson user,
                                    PublicKeyJson cosigningPublicKey,
                                    PublicKeyJson swapServerPublicKey) {

        this.user = user;
        this.cosigningPublicKey = cosigningPublicKey;
        this.swapServerPublicKey = swapServerPublicKey;
    }
}

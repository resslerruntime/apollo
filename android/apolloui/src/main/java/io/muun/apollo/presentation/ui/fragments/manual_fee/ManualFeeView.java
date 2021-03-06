package io.muun.apollo.presentation.ui.fragments.manual_fee;

import io.muun.apollo.domain.model.CurrencyDisplayMode;
import io.muun.apollo.domain.model.PaymentContext;
import io.muun.apollo.domain.model.PaymentRequest;
import io.muun.apollo.presentation.ui.base.BaseView;

interface ManualFeeView extends BaseView {

    void setCurrencyDisplayMode(CurrencyDisplayMode currencyDisplayMode);

    void setPaymentContext(PaymentContext paymentContext, PaymentRequest paymentRequest);
}

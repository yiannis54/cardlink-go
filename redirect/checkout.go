package redirect

func (s *Signer) CheckoutURL(reqParams *PaymentRequest) (string, error) {
	u, err := s.Config.RedirectURL()
	if err != nil {
		return "", err
	}

	fields, err := s.Sign(reqParams)
	if err != nil {
		return "", err
	}

	return FormHTML(u.String(), fields), nil
}

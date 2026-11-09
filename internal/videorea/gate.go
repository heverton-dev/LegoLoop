package videorea

// EvaluateStopGate valida a decisão do operador em qualquer stop-gate do pipeline.
func EvaluateStopGate(gate StopGateType, approved bool, notes string) (*StopGateDecision, error) {
	d := &StopGateDecision{
		Gate:     gate,
		Approved: approved,
		Notes:    notes,
	}
	if !approved {
		return d, ErrGateRejected
	}
	return d, nil
}

// EvaluateCurationGate valida a decisão do operador no stop-gate de curadoria.
func EvaluateCurationGate(report *CurationReport, approved bool) error {
	if report == nil {
		return ErrEmptyDossier
	}
	_, err := EvaluateStopGate(GateCuradoria, approved, "")
	if err != nil {
		report.Approved = false
		return err
	}
	report.Approved = true
	return nil
}

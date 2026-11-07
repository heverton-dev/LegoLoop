package videorea

// EvaluateCurationGate valida a decisão do operador no stop-gate de curadoria.
func EvaluateCurationGate(report *CurationReport, approved bool) error {
	if report == nil {
		return ErrEmptyDossier
	}
	if !approved {
		report.Approved = false
		return ErrGateRejected
	}
	report.Approved = true
	return nil
}

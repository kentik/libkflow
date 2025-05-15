package flow

// NormalizeSampleRate adjusts the sample rate in place on the provided [flow.Flow] slice based on a provided
// adjustment factor if it is > 1.0. The adjustment factor is multiplied by the original sample rate and 100 to get
// the new sample rate, as it is expected that a [flow.Flow] with a sample rate that does not account for this change.
func NormalizeSampleRate(flows []Flow, resampleRateAdj float32) {
	for i := range flows {
		sampleRate := flows[i].SampleRate
		adjustedSR := sampleRate * 100

		if resampleRateAdj > 1.0 {
			adjustedSR = uint32(float32(adjustedSR) * resampleRateAdj)
		}

		flows[i].SampleAdj = true
		flows[i].SampleRate = adjustedSR
	}
}

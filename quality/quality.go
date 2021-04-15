package quality

import (
	"errors"
)

const (
	defaultPhred = 33
	defaultStrScore = 'H'
	defaultIntScore = 40
)

type Quality struct {
	Phred    int
	IntScore []int
	StrScore []byte
}

// Quality builder
func NewQuality(phred int) *Quality {
	q := Quality{Phred: phred}
	return &q
}

func (q *Quality) identifyPhredFromBytes(score []byte) error {
	min := int(score[0])
	max := int(score[0])
	for i:=1 ; i<len(score) ; i++ {
		val := int(score[i])
		if min > val {
			min = val
		} else {
			if max < val {
				max = val
			}
		}
	}

	// Decide which Phred type
	if min < 64 {
		if max < 74 && min > 32 {
			// This is PHRED 33 score
			q.Phred = 33
			return nil
		} else {
			return errors.New("[PRHED QUALITY]: Unconsistant PHRED values.")
		}
	}
	if max > 73 {
		if min > 63 && max < 105 {
			// This is PHRED 64
			q.Phred = 64
			return nil
		} else {
			return errors.New("[PHRED QUALITY]: Unconsistant PHRED values.")
		}
	}

	// NOTE: PHRED type not determined
	return errors.New("[PHRED QUALITY]: Failed to determine the PHRED type.")
}

func (q *Quality) appendIntScoreFromByte(score []byte) {
	phred := defaultPhred
	if q.Phred != 0 {
		phred = q.Phred
	}
	for _, b := range score {
		q.IntScore = append(q.IntScore, int(b)- phred)
	}
}

// Append Phred score from bytes (srt)
func (q *Quality) AppendStrScore(score []byte) error {
	var err error

	// Uninitialized Phred type
	if q.Phred == 0 {
		err = q.identifyPhredFromBytes(score)
	}

	// Append the quality score (str)
	q.StrScore = append(q.StrScore, score...)

	// Append the quality score (int)
	q.appendIntScoreFromByte(score)

	return err
}

// Generate a default score entry
func (q *Quality) GenerateDefaultScore(n int) {
	q.Phred = defaultPhred
	for i := 0 ; i<n ; i++ {
		q.StrScore = append(q.StrScore, defaultStrScore)
		q.IntScore = append(q.IntScore, defaultIntScore)
	}
}
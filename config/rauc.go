package config

const RaucBundle Bundle = "rauc"

// Rauc configures an A/B system update using RAUC.
//
// example payload
// {
//   "pre_condition": "true",
//   "rauc_bundle": "/path/to/bundle.raucb",
// }

type Rauc struct {
	Metadata

	// PreCondition defines an optional command which needs to return 0 in order for RAUC bundle to be installed.
	PreCondition string `json:"pre_condition,omitempty"`

	// RaucBundle defines the rauc bundle to be installed.
	RaucBundle string `json:"rauc_bundle"`
}

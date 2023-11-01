package config

// Pushed is a struct that represents the pushed configuration.
type Pushed struct {
	// CommitID represents commit ID of the mose recent commit affecting the device.
	CommitID string `json:"commit_id"`

	// Bundles contains a list of strings representing configuration bundles
	Bundles BundleNames `json:"bundles"`

	// BundleData contain configuration data for bundles in the Bundles list
	BundleData BundleData `json:"bundle_data"`
}

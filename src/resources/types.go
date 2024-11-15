package resources

type ResourceMap map[string]ResourceSpec

type ScalarMessage struct {
	Value string
}

type VectorMessage struct {
	Values []string
}

type Resources struct {
	registeredResources ResourceMap
}

type ResourceSpec struct {
	URL string					`json:"url"`
	QuerySelector string		`json:"querySelector"`
	QuerySelectorAll string		`json:"querySelectorAll"`
}

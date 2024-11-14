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
	URL string
	QuerySelector string
	QuerySelectorAll string
}

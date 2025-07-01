package mongomodels


type EnvVar struct {
	Name  string `bson:"name" json:"name"`
	Value string `bson:"value" json:"value"`
}

type ContainerConfig struct {
	Envs []EnvVar `bson:"envs" json:"envs"`
}

type WorkloadConfig struct {
	ID         string                       `bson:"_id" json:"_id"` // MongoDB document ID
	Name       string                       `bson:"name" json:"name"`
	Containers map[string]ContainerConfig   `bson:"containers" json:"containers"`
}

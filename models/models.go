package models

type AdmissionReview struct {
	Request Request `json:"request"`
}

type Request struct {
	UID       string     `json:"uid"`
	UserInfo  UserInfo   `json:"userInfo"`
	Name      string     `json:"name"`
	Namespace string     `json:"namespace"`
	Operation string     `json:"operation"`
	Object    ObjectInfo `json:"object"`
}

type UserInfo struct {
	Username string   `json:"username"`
	UID      string   `json:"uid"`
	Groups   []string `json:"groups"`
}

type ObjectInfo struct {
	Metadata Metadata `json:"metadata"`
	Spec     Spec     `json:"spec"`
}

type Metadata struct {
	UID string `json:"uid"`
}

type Spec struct {
	Template Template `json:"template"`
}

type Template struct {
	Spec PodSpec `json:"spec"`
}

type PodSpec struct {
	Containers []Container `json:"containers"`
}

type Container struct {
	Name            string `json:"name"`
	Image           string `json:"image"`
	ImagePullPolicy string `json:"imagePullPolicy"`
}


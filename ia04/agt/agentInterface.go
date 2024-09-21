package agt

// AgentI defines the methods that each agent must implement.
type AgentI interface {
	Equal(ag AgentI) bool
	DeepEqual(ag AgentI) bool
	Clone() AgentI
	String() string
	Prefers(a Alternative, b Alternative) bool
	Start()
}

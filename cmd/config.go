package cmd

type Config struct {
	Simulation  Simulation
	TaskDir     string
	DoneDir     string
	Repository  []Repository
	SpreadSheet SpreadSheet
}

type Simulation struct {
	Monte  []string
	Range  Range
	SimDir string
	DstDir string
	Vtp    Node
	Vtn    Node
}

type SpreadSheet struct {
	Id        string
	CSPath    string
	TokenPath string
}

type Task struct {
	Simulation Simulation
}

type RepoType int

const (
	Git RepoType = iota
	AWS RepoType = iota
)

type Repository struct {
	Type        RepoType
	Path        string
	DirPattern  string
	FilePattern string
}

type Node struct {
	Voltage   float64
	Sigma     float64
	Deviation float64
}

type Range struct {
	Start string
	Stop  string
	Step  string
}

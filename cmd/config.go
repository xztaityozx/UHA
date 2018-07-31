package cmd

type Config struct {
	Simulation  Simulation
	TaskDir     string
	Repository  []Repository
	SpreadSheet SpreadSheet
	SlackConfig SlackConfig
}

var ConfigDir string
var ReserveRunDir string
var DoneRunDir string
var FailedRunDir string
var ReserveSRunDir string
var DoneSRunDir string
var FailedSRunDir string
var NextPath string

type Simulation struct {
	Monte  []string
	Range  Range
	Signal string
	SimDir string
	DstDir string
	Vtp    Node
	Vtn    Node
}

type SpreadSheet struct {
	Id        string
	CSPath    string
	TokenPath string
	SheetName string
}

type Task struct {
	Simulation Simulation
}

const (
	Git string = "Git"
	AWS string = "AWS"
	Dir string = "Dir"
)

type Repository struct {
	Type string
	Path string
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

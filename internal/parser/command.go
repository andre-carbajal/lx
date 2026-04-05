package parser

type ParsedCommand struct {
	Raw   string         // Input original: "ls -la /home"
	Name  string         // Nombre del comando: "ls"
	Flags []string       // Flags individuales: ["-l", "-a"]
	Args  []string       // Argumentos: ["/home"]
	Pipe  *ParsedCommand // Siguiente comando si hay pipe (nil si no)
}

package essh

func ExamplePrintVersion() {
	Run([]string{"--version"})
	// Output:
	// dev (unknown)
}
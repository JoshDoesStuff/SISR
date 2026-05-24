package steam

func ExecuteableDir() (string, error) {
	return steamPath()
}

func ClientRunning() bool {
	return steamRunning()
}

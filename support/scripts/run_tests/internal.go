package main

type packageInfo struct {
}

func loadPackageInfo() ([]packageInfo, error) {
	cmd := exec.Command("go", "list", "-json")
	return nil, nil
}

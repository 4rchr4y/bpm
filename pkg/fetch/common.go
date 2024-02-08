package fetch

// func shouldIgnore(ignoreList bundle.IgnoreList, path string) bool {
// 	if path == "" || len(ignoreList) == 0 {
// 		return false
// 	}

// 	dir := filepath.Dir(path)
// 	if dir == "." {
// 		return false
// 	}

// 	topLevelDir := strings.Split(dir, string(filepath.Separator))[0]
// 	_, found := ignoreList[topLevelDir]
// 	return found
// }

// func readLinesToMap(content []byte) (bundle.IgnoreList, error) {
// 	result := make(bundle.IgnoreList)
// 	scanner := bufio.NewScanner(bytes.NewReader(content))
// 	for scanner.Scan() {
// 		line := strings.TrimSpace(scanner.Text())
// 		if line == "" {
// 			continue
// 		}

// 		result[line] = struct{}{}
// 	}

// 	if err := scanner.Err(); err != nil {
// 		return nil, fmt.Errorf("error reading '%s' input: %v", constant.IgnoreFileName, err)
// 	}

// 	return result, nil
// }

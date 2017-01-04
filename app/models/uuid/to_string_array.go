package uuid

func ToStringArray(uuids []UUID) []string {
	stringArray := make([]string, len(uuids))
	for i, uuid := range uuids {
		stringArray[i] = uuid.String()
	}
	return stringArray
}

package main

import "transcoder/internal/converter"

func main() {
	converter := converter.NewVideoConvert()
	converter.Handle([]byte(`{"video_id": 1, "path": "/media/uploads/1"}`))
}

package main

type deleteFileResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

const deleteFileURL = "https://slack.com/api/files.delete?token=%v&file=%v"

type fileResponse struct {
	ID string `json:"id"`
}

type listFilesResponse struct {
	Files []fileResponse `json:"files"`
}

const listFilesURL = "https://slack.com/api/files.list?token=%v&ts_to=%v"

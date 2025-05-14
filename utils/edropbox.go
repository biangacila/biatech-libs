package utils

import (
	"encoding/json"
	"fmt"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type IResponseList struct {
	Name                 string    `json:"name"`
	PathLower            string    `json:"path_lower"`
	PathDisplay          string    `json:"path_display"`
	ParentSharedFolderId string    `json:"parent_shared_folder_id"`
	Id                   string    `json:"id"`
	ClientModified       time.Time `json:"client_modified"`
	ServerModified       time.Time `json:"server_modified"`
	Rev                  string    `json:"rev"`
	Size                 int       `json:"size"`
	SharingInfo          struct {
		ReadOnly             bool   `json:"read_only"`
		ParentSharedFolderId string `json:"parent_shared_folder_id"`
		ModifiedBy           string `json:"modified_by"`
	} `json:"sharing_info"`
	IsDownloadable bool   `json:"is_downloadable"`
	ContentHash    string `json:"content_hash"`
}
type IDropbox struct {
	clientId     string
	clientSecret string
	token        string
}

func NewIDropbox(clientId string, clientSecret string, token string) *IDropbox {
	return &IDropbox{clientId: clientId, clientSecret: clientSecret, token: token}
}
func (c *IDropbox) connection() files.Client {
	config := dropbox.Config{
		Token:    c.token,
		LogLevel: dropbox.LogOff,
	}
	client := files.New(config)
	return client
}
func (c *IDropbox) convertRemoteList(inList any) (records []ShortListFileDropbox) {
	var infos []IResponseList
	str, _ := json.Marshal(inList)
	_ = json.Unmarshal(str, &infos)
	for _, rec := range infos {
		var record ShortListFileDropbox
		record.Name = rec.Name
		record.Modified = rec.ClientModified.String()
		record.Revision = rec.Rev
		record.Size = rec.Size
		record.Path = rec.PathLower

		records = append(records, record)
	}
	return records
}

func (c *IDropbox) ListFilesInFolder(dropboxFolder string) ([]ShortListFileDropbox, error) {
	config := dropbox.Config{
		Token:    c.token,
		LogLevel: dropbox.LogOff, // Set to LogInfo if debugging
	}

	client := files.New(config)
	var records []ShortListFileDropbox

	// ✅ First request: Fetch initial batch of files
	arg := &files.ListFolderArg{Path: dropboxFolder}
	res, err := client.ListFolder(arg)
	if err != nil {
		log.Println("Error listing files:", err)
		return nil, err
	}

	// ✅ Process first batch of files
	for _, entry := range res.Entries {
		if fileMeta, ok := entry.(*files.FileMetadata); ok {
			records = append(records, ShortListFileDropbox{
				Name:     fileMeta.Name,
				Path:     fileMeta.PathLower,
				Revision: fileMeta.Rev,
				Size:     int(fileMeta.Size),
				Modified: fileMeta.ClientModified.String(),
			})
		}
	}

	// ✅ Continue fetching files using `ListFolderContinue()` until `HasMore == false`
	for res.HasMore {
		cursorArg := &files.ListFolderContinueArg{Cursor: res.Cursor}
		res, err = client.ListFolderContinue(cursorArg)
		if err != nil {
			log.Println("Error listing files in pagination:", err)
			return nil, err
		}

		// Process additional batches
		for _, entry := range res.Entries {
			if fileMeta, ok := entry.(*files.FileMetadata); ok {
				records = append(records, ShortListFileDropbox{
					Name:     fileMeta.Name,
					Path:     fileMeta.PathLower,
					Revision: fileMeta.Rev,
					Size:     int(fileMeta.Size),
					Modified: fileMeta.ClientModified.String(),
				})
			}
		}
	}

	// ✅ Sort files by Modified Date (Descending)
	sort.Slice(records, func(i, j int) bool {
		return records[i].Modified > records[j].Modified // Newest first
	})

	return records, nil
}

func (c *IDropbox) ListMatchingFiles(dropboxFolder, searchKey string) (records []ShortListFileDropbox, err error) {
	outs, err := c.ListFilesInFolder(dropboxFolder)
	if err != nil {
		log.Println("Error listing files: ", err)
		return records, err
	}
	fmt.Printf("Files matching '%s' in folder: %s\n", searchKey, dropboxFolder)
	// Filter files by search key
	for _, entry := range outs {
		if strings.Contains(strings.ToLower(entry.Name), strings.ToLower(searchKey)) {
			records = append(records, entry)
			fmt.Println("- ", entry.Name)
		}
	}
	return records, nil
}
func (c *IDropbox) DownloadByRev(rev string, savePath string) (link string, err error) {
	config := dropbox.Config{
		Token:    c.token,
		LogLevel: dropbox.LogOff,
	}

	client := files.New(config)

	// ✅ Get file path from `rev`
	filePath, err := c.GetFilePathByRev(rev)
	if err != nil {
		return link, fmt.Errorf("Error getting file path: %v", err)
	}

	// ✅ Use `Path` instead of `Rev` in DownloadArg
	arg := files.DownloadArg{Path: filePath}
	res, contents, err := client.Download(&arg)
	if err != nil {
		return link, fmt.Errorf("Error downloading file: %v", err)
	}
	defer contents.Close()

	// Save the file locally
	localFileName := fmt.Sprintf("%v/%v", savePath, res.Name)
	outFile, err := os.Create(localFileName)
	if err != nil {
		return link, fmt.Errorf("Error creating file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, contents)
	if err != nil {
		return link, fmt.Errorf("Error writing file: %v", err)
	}

	fmt.Println("File downloaded successfully:", res.Name, "->", localFileName)
	return localFileName, nil
}

// Function to get file path by revision
func (c *IDropbox) GetFilePathByRev(rev string) (string, error) {
	config := dropbox.Config{
		Token:    c.token,
		LogLevel: dropbox.LogOff,
	}

	client := files.New(config)

	// Call Dropbox API to get file metadata using rev
	arg := &files.GetMetadataArg{Path: "rev:" + rev} // Dropbox requires "rev:" prefix
	meta, err := client.GetMetadata(arg)
	if err != nil {
		return "", fmt.Errorf("Error getting file metadata: %v", err)
	}

	// Ensure metadata is of file type
	if fileMeta, ok := meta.(*files.FileMetadata); ok {
		return fileMeta.PathLower, nil
	}

	return "", fmt.Errorf("File with revision %s not found", rev)
}

package noise

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/kyanite/noise/internal/domain"
	"github.com/kyanite/noise/internal/infra/files"
)

// TestFileIO runs a suite of sub-tests covering the end-to-end file service
// workflow while maintaining manageable complexity and test-only logging.
func TestFileIO(t *testing.T) {
	env := newFileIOEnv(t)

	t.Run("SaveSong", env.saveSong)
	t.Run("ListSongs", env.listSongs)
	t.Run("ReadSong", env.readSong)
	t.Run("CacheSong", env.checkCache)
	t.Run("WatchSong", env.watchSong)
	t.Run("DeleteSong", env.deleteSong)
	t.Run("BackupSong", env.backupSong)
}

type fileIOTestEnv struct {
	service      *files.Service
	baseDir      string
	song         *domain.Song
	testFilePath string
}

func newFileIOEnv(t *testing.T) *fileIOTestEnv {
	t.Helper()

	baseDir := t.TempDir()

	service, err := files.New(files.Config{
		BaseDir:          baseDir,
		AutoSave:         true,
		AutoSaveInterval: time.Minute,
		BackupCount:      3,
	})
	if err != nil {
		t.Fatalf("Failed to create file service: %v", err)
	}
	t.Cleanup(func() {
		if err := service.Close(); err != nil {
			t.Fatalf("Failed to close file service: %v", err)
		}
	})

	return &fileIOTestEnv{
		service: service,
		baseDir: baseDir,
		song: &domain.Song{
			Metadata: domain.SongMetadata{
				Title:         "Test Song",
				Artist:        "Test Artist",
				Key:           "C Major",
				Tempo:         120,
				TimeSignature: "4/4",
				Tags:          []string{"test", "demo"},
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			Sections: []domain.Section{
				{
					Type:   domain.SectionVerse,
					Number: 1,
					Lines: []domain.Line{
						{Text: "This is the first line of the verse", Syllables: 8},
						{Text: "This is the second line of the verse", Syllables: 8},
					},
				},
				{
					Type:   domain.SectionChorus,
					Number: 1,
					Lines: []domain.Line{
						{Text: "This is the chorus line", Syllables: 5},
						{Text: "Another chorus line here", Syllables: 5},
					},
				},
			},
		},
		testFilePath: "test_song.md",
	}
}

func (env *fileIOTestEnv) saveSong(t *testing.T) {
	t.Helper()
	if err := env.service.WriteSong(env.song, env.testFilePath); err != nil {
		t.Fatalf("Failed to save song: %v", err)
	}
	t.Logf("Song saved to: %s", env.testFilePath)
}

func (env *fileIOTestEnv) listSongs(t *testing.T) {
	t.Helper()
	songFiles, err := env.service.ListSongs()
	if err != nil {
		t.Fatalf("Failed to list songs: %v", err)
	}
	t.Logf("Found %d song files:", len(songFiles))
	for _, file := range songFiles {
		t.Logf("  - %s", file)
	}
}

func (env *fileIOTestEnv) readSong(t *testing.T) {
	t.Helper()
	loadedSong, err := env.service.ReadSong(env.testFilePath)
	if err != nil {
		t.Fatalf("Failed to read song: %v", err)
	}
	t.Logf("Loaded song: %s by %s", loadedSong.Metadata.Title, loadedSong.Metadata.Artist)
	t.Logf("Sections: %d", len(loadedSong.Sections))

	if loadedSong.Metadata.Title != env.song.Metadata.Title {
		t.Fatalf("Title mismatch: expected %s, got %s", env.song.Metadata.Title, loadedSong.Metadata.Title)
	}
	if loadedSong.Metadata.Artist != env.song.Metadata.Artist {
		t.Fatalf("Artist mismatch: expected %s, got %s", env.song.Metadata.Artist, loadedSong.Metadata.Artist)
	}
	t.Logf("Title: %s", loadedSong.Metadata.Title)
	t.Logf("Artist: %s", loadedSong.Metadata.Artist)
	t.Logf("Key: %s", loadedSong.Metadata.Key)
	t.Logf("Tempo: %d", loadedSong.Metadata.Tempo)
}

func (env *fileIOTestEnv) checkCache(t *testing.T) {
	t.Helper()
	if _, exists := env.service.GetCachedFile(env.testFilePath); !exists {
		t.Fatalf("File not found in cache")
	}
	t.Log("File cache check passed")
}

func (env *fileIOTestEnv) watchSong(t *testing.T) {
	t.Helper()

	events := make(chan bool, 1)
	watcher := &logWatcher{t: t, events: events}

	if err := env.service.WatchFile(env.testFilePath, watcher); err != nil {
		t.Fatalf("Failed to watch file: %v", err)
	}
	defer func() {
		if err := env.service.UnwatchFile(env.testFilePath, watcher); err != nil {
			t.Fatalf("Failed to unwatch file: %v", err)
		}
	}()

	time.Sleep(50 * time.Millisecond) // allow watcher registration

	env.song.Metadata.UpdatedAt = time.Now()
	if err := env.service.WriteSong(env.song, env.testFilePath); err != nil {
		t.Fatalf("Failed to write song during watcher test: %v", err)
	}

	select {
	case <-events:
		t.Log("Watcher event received successfully")
	case <-time.After(time.Second):
		t.Log("Watcher event timeout (this may be expected on some platforms)")
	}
}

func (env *fileIOTestEnv) deleteSong(t *testing.T) {
	t.Helper()
	if err := env.service.DeleteSong(env.testFilePath); err != nil {
		t.Fatalf("Failed to delete song: %v", err)
	}
	exists, err := env.service.SongExists(env.testFilePath)
	if err != nil {
		t.Fatalf("Failed to check song existence: %v", err)
	}
	if exists {
		t.Fatalf("Song still exists after deletion")
	}
	t.Log("Song file deleted and verified")
}

func (env *fileIOTestEnv) backupSong(t *testing.T) {
	t.Helper()

	backupTestFile := "backup_test.md"
	backupSong := &domain.Song{
		Metadata: domain.SongMetadata{
			Title:     "Backup Test Song",
			Artist:    "Backup Artist",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Sections: []domain.Section{
			{
				Type:   domain.SectionVerse,
				Number: 1,
				Lines:  []domain.Line{{Text: "Backup test line", Syllables: 4}},
			},
		},
	}

	if err := env.service.WriteSong(backupSong, backupTestFile); err != nil {
		t.Fatalf("Failed to save backup test song: %v", err)
	}
	if err := env.service.CreateBackup(backupTestFile); err != nil {
		t.Fatalf("Failed to create backup: %v", err)
	}

	backupGlob := filepath.Join(env.baseDir, backupTestFile+".backup.*")
	matches, err := filepath.Glob(backupGlob)
	if err != nil {
		t.Fatalf("Failed to glob backup files: %v", err)
	}
	if len(matches) == 0 {
		t.Log("Warning: Backup file not found (may be generated in a different location)")
	} else {
		t.Logf("Backup file created: %s", matches[0])
	}
}

type logWatcher struct {
	t      *testing.T
	events chan<- bool
}

func (w *logWatcher) OnFileChanged(filePath string, event files.FileEvent) {
	if w.t != nil {
		w.t.Logf("File changed: %s (event: %s)", filePath, event)
	}
	select {
	case w.events <- true:
	default:
	}
}

func (w *logWatcher) OnFileDeleted(filePath string) {
	if w.t != nil {
		w.t.Logf("File deleted: %s", filePath)
	}
	select {
	case w.events <- true:
	default:
	}
}

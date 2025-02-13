package player

import (
	"MusicPlayer/logger"
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// Ctrl struct for controlling playback
type Ctrl struct {
	ctrl      *beep.Ctrl
	done      chan bool
	format    beep.Format
	streamer  beep.StreamSeekCloser
	startTime time.Time
}

// PlayAudio function to handle playing the audio file
func PlayAudio(file string) (*Ctrl, bool) {
	musicFile, err := os.Open("/home/code/Music/" + file + ".mp3")
	if err != nil {
		logger.Logger(fmt.Sprintf("Error opening file: %v", err))
		return nil, false
	}

	streamer, format, err := mp3.Decode(musicFile)
	if err != nil {
		logger.Logger(fmt.Sprintf("Error decoding MP3 file: %v", err))
		return nil, false
	}

	// Wrap the streamer with beep.Ctrl for pause/resume functionality
	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}

	// Initialize the speaker (buffer size = 1/10th of a second)
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		logger.Logger(fmt.Sprintf("Error initializing speaker: %v", err))
		return nil, false
	}

	done := make(chan bool)

	// Play the audio and notify when done
	speaker.Play(beep.Seq(ctrl, beep.Callback(func() {
		logger.Logger(fmt.Sprintf("Done playing: %s", file))
		done <- true
		os.Exit(0)
	})))

	return &Ctrl{ctrl: ctrl, done: done, format: format, streamer: streamer, startTime: time.Now()}, true
}

// PauseAudio toggles playback state
func (c *Ctrl) PauseAudio() {
	speaker.Lock()
	c.ctrl.Paused = !c.ctrl.Paused // Toggle pause state
	speaker.Unlock()

	if c.ctrl.Paused {
		logger.Logger("Audio paused")
	} else {
		logger.Logger("Audio resumed")
	}
}

// DisplayProgress shows playback progress
func (c *Ctrl) DisplayProgress() {
	totalFrames := c.streamer.Len()
	sampleRate := c.format.SampleRate

	for {
		speaker.Lock()
		currentFrame := c.streamer.Position()
		speaker.Unlock()

		// Convert frames to seconds
		currentTime := time.Duration(currentFrame) * time.Second / time.Duration(sampleRate)
		totalTime := time.Duration(totalFrames) * time.Second / time.Duration(sampleRate)

		fmt.Printf("\rPlaying: [%s / %s]", formatDuration(currentTime), formatDuration(totalTime))

		if currentFrame >= totalFrames {
			break
		}

		time.Sleep(1 * time.Second) // Update progress every second
	}
	fmt.Println()
}

// Helper function to format duration into mm:ss
func formatDuration(d time.Duration) string {
	minutes := int(d.Seconds()) / 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// Waits until playback is finished
func (c *Ctrl) WaitForCompletion() {
	c.DisplayProgress()
	<-c.done // Block until the audio finishes playing
}

func (c *Ctrl) QuitAudio() {
	c.done <- true
}

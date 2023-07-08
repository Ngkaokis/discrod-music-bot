package types

// ActiveGuild is a guild that is currently being streamed to.
type ActiveGuild struct {
	Name        string
	MediaChan   chan *Media
	UserActions *UserActions
}

func NewActiveGuild(name string) *ActiveGuild {
	return &ActiveGuild{
		Name:        name,
		MediaChan:   nil,
		UserActions: nil,
	}
}

func (g *ActiveGuild) PrepareForStreaming(maxQueueSize int) {
	g.MediaChan = make(chan *Media, maxQueueSize)
	g.UserActions = NewActions()
}

func (g *ActiveGuild) IsStreaming() bool {
	return g.MediaChan != nil
}

func (g *ActiveGuild) EnqueueMedia(media *Media) {
	g.MediaChan <- media
}

func (g *ActiveGuild) MediaQueueSize() int {
	return len(g.MediaChan)
}

func (g *ActiveGuild) IsMediaQueueFull() bool {
	return g.MediaChan != nil && len(g.MediaChan) == cap(g.MediaChan)
}

func (g *ActiveGuild) StopStreaming() {
	close(g.MediaChan)
	g.MediaChan = nil
	g.UserActions = nil
}
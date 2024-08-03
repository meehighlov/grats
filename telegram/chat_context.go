package telegram

type ChatContext interface {
	// appends text to context
	AppendText(string) error

	// gets all saved text in context
	GetTexts() []string

	getCommandInProgress() string
	setCommandInProgress(string) error
	getStepTODO() int
	setStepTODO(int) error
	reset() error
}

type chatContext struct {
	chatId            string
	userResponses     []string
	commandInProgress string
	stepTODO          int
}

func newChatContext(chatId string) *chatContext {
	return &chatContext{chatId, []string{}, "", ENTRYPOINT_STEP}
}

func (ctx *chatContext) AppendText(userResponse string) error {
	ctx.userResponses = append(ctx.userResponses, userResponse)
	return nil
}

func (ctx *chatContext) GetTexts() []string {
	return ctx.userResponses
}

func (ctx *chatContext) getCommandInProgress() string {
	return ctx.commandInProgress
}

func (ctx *chatContext) getStepTODO() int {
	return ctx.stepTODO
}

func (ctx *chatContext) setStepTODO(step int) error {
	ctx.stepTODO = step
	return nil
}

func (ctx *chatContext) setCommandInProgress(command string) error {
	if ctx.commandInProgress != "" {
		if ctx.commandInProgress != command {
			ctx.reset()
		}
	}
	ctx.commandInProgress = command
	return nil
}

func (ctx *chatContext) reset() error {
	ctx.commandInProgress = ""
	ctx.userResponses = []string{}
	ctx.stepTODO = ENTRYPOINT_STEP
	return nil
}

func (bot *bot) getOrCreateChatContext(chatId string) ChatContext {
	ctx, found := bot.cache.Get(chatId)

	if found {
		return ctx.(*chatContext)
	}

	newCtx := newChatContext(chatId)

	bot.cache.Set(chatId, newCtx, bot.cacheExparation)

	return newCtx
}

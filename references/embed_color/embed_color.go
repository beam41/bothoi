package embed_color

type EmbedColor uint32

const (
	Default              EmbedColor = 0xb8b8b8
	Info                 EmbedColor = 0x0d47a1
	Error                EmbedColor = 0xd32f2f
	ErrorLow             EmbedColor = 0xf50057
	SuccessInterrupt     EmbedColor = 0xeeff41
	SuccessInterruptHigh EmbedColor = 0xff6d00
	SuccessContinue      EmbedColor = 0x9ccc65
	SuccessScheduled     EmbedColor = 0xb2ebf2
)

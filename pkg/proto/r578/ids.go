package r578

// Handshaking State
// Clientbound packet IDs
const (
	HandshakingHandshake int = iota // 0x00
)

// Login State
// Clientbound packet IDs
const (
	LoginDisconnect        int = iota // 0x00
	LoginEncryptionRequest            // 0x01
	LoginSuccess                      // 0x02
	LoginSetCompression               // 0x03
	LoginPluginRequest                // 0x04
)

// Serverbound packet IDs
const (
	LoginStart              int32 = iota // 0x00
	LoginEncryptionResponse              // 0x01
	LoginPluginResponse                  // 0x02
)

// Play State
// Clientbound packet IDs
const (
	PlaySpawnEntity        int = iota // 0x00
	PlaySpawnExperienceOrb            // 0x01
	PlaySpawnWeatherEntity            // 0x02
	PlaySpawnLivingEntity             // 0x03
	PlaySpawnPainting                 // 0x04
	PlaySpawnPlayer                   // 0x05
	PlayEntityAnimation               // 0x06
	PlayStatistics                    // 0x07

)

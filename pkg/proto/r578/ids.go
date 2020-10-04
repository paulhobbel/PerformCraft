package r578

import "github.com/paulhobbel/performcraft/pkg/common"

// Handshake State
// Clientbound packet IDs
const (
	HandshakingHandshake common.PacketID = iota // 0x00
)

// Status State
// Clientbound packet IDs
const (
	StatusRequest common.PacketID = iota // 0x00
	StatusPing                           // 0x01
)

const (
	StatusResponse common.PacketID = iota // 0x00
	StatusPong                            // 0x01
)

// Login State
// Clientbound packet IDs
const (
	LoginDisconnect        common.PacketID = iota // 0x00
	LoginEncryptionRequest                        // 0x01
	LoginSuccess                                  // 0x02
	LoginSetCompression                           // 0x03
	LoginPluginRequest                            // 0x04
)

// Serverbound packet IDs
const (
	LoginStart              common.PacketID = iota // 0x00
	LoginEncryptionResponse                        // 0x01
	LoginPluginResponse                            // 0x02
)

// Play State
// Clientbound packet IDs
const ()

// Serverbound packetIDs
const (
	PlaySpawnEntity              common.PacketID = iota // 0x00
	PlaySpawnExperienceOrb                              // 0x01
	PlaySpawnWeatherEntity                              // 0x02
	PlaySpawnLivingEntity                               // 0x03
	PlaySpawnPainting                                   // 0x04
	PlaySpawnPlayer                                     // 0x05
	PlayEntityAnimation                                 // 0x06
	PlayStatistics                                      // 0x07
	PlayAckPlayerDig                                    // 0x08
	PlayBlockBreakAnimation                             // 0x09
	PlayBlockEntityData                                 // 0x0A
	PlayBlockAction                                     // 0x0B
	PlayBlockChange                                     // 0x0C
	PlayBossBar                                         // 0x0D
	PlayServerDifficulty                                // 0x0E
	PlayChatMessageServer                               // 0x0F
	PlayMultiBlockChange                                // 0x10
	PlayTabCompleteServer                               // 0x11
	PlayDeclareCommands                                 // 0x12
	PlayWindowConfirmServer                             // 0x13
	PlayCloseWindowServer                               // 0x14
	PlayWindowItems                                     // 0x15
	PlayWindowProperty                                  // 0x16
	PlaySetSlot                                         // 0x17
	PlaySetCooldown                                     // 0x18
	PlayPluginMessageServer                             // 0x19
	PlayNamedSoundEffect                                // 0x1A
	PlayDisconnect                                      // 0x1B
	PlayEntityStatus                                    // 0x1C
	PlayExplosion                                       // 0x1D
	PlayUnloadChuck                                     // 0x1E
	PlayChangeGameState                                 // 0x1F
	PlayOpenHorseWindow                                 // 0x20
	PlayKeepAliveServer                                 // 0x21
	PlayChuckData                                       // 0x22
	PlayEffect                                          // 0x23
	PlayParticle                                        // 0x24
	PlayUpdateLight                                     // 0x25
	PlayJoinGame                                        // 0x26
	PlayMapData                                         // 0x27
	PlayTradeList                                       // 0x28
	PlayEntityPosition                                  // 0x29
	PlayEntityPositionRotation                          // 0x2A
	PlayEntityRotation                                  // 0x2B
	PlayEntityMovement                                  // 0x2C
	PlayVehicleMoveServer                               // 0x2D
	PlayOpenBook                                        // 0x2E
	PlayOpenWindow                                      // 0x2F
	PlayOpenSignEditor                                  // 0x30
	PlayCraftRecipeResponse                             // 0x31
	PlayPlayerAbilitiesServer                           // 0x32
	PlayCombatEvent                                     // 0x33
	PlayPlayerInfo                                      // 0x34
	PlayFacePlayer                                      // 0x35
	PlayPlayerPositionLookServer                        // 0x36
)

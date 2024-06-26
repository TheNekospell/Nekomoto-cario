package database

import (
	// "backend/internal/chain"
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func GetNekoSpiritInfoByTokenId(tokenId uint64) (*ServerNekoSpiritInfo, error) {
	var boxInfo ServerNekoSpiritInfo
	if err := DB.Where("token_id = ?", tokenId).First(&boxInfo).Error; err != nil {
		return nil, err
	}
	return &boxInfo, nil
}

func GetUidByAddress(owner common.Address) (uint64, error) {
	var user ServerAddress
	if err := DB.Where("address = ?", owner.Hex()).First(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func GetUidByInviteCode(inviteCode string) (uint64, error) {
	var user ServerAddress
	if err := DB.Where("invite_code = ?", inviteCode).First(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func CreateNekoSpiritInfo(box *ServerNekoSpiritInfo) error {
	if err := DB.Create(&box).Error; err != nil {
		return err
	}
	Cache.Delete(CacheTagUid + strconv.FormatUint(box.StakeFromUid, 10))
	return nil
}

func UpdateNekoSpiritInfo(box *ServerNekoSpiritInfo) error {
	if err := DB.Model(&box).Updates(box).Error; err != nil {
		return err
	}
	Cache.Delete(CacheTagNekoSpirit + strconv.FormatUint(box.TokenId, 10))
	return nil
}

func CreateBuffRecord(record *ServerBuffRecord) error {
	if err := DB.Create(&record).Error; err != nil {
		return err
	}
	return nil
}

func UpdateBuffRecord(record *ServerBuffRecord) error {
	if err := DB.Model(&record).Updates(record).Error; err != nil {
		return err
	}
	Cache.Delete(CacheTagUid + strconv.FormatUint(record.Uid, 10))
	return nil
}

// func SaveEventNekoSpiritTransfer(event *chain.SpiritTransfer) error {
// 	toSave := EventNekoSpiritTransfer{
// 		Hash: Hash{
// 			TransactionHash: event.Raw.TxHash.Hex(),
// 			BlockNumber:     event.Raw.BlockNumber,
// 			BlockHash:       event.Raw.BlockHash.Hex()},
// 		From:    event.From.Hex(),
// 		To:      event.To.Hex(),
// 		TokenId: event.TokenId.Uint64(),
// 	}

// 	if toSave.From == chain.EmptyAddressInHex {
// 		toSave.TransferOrMintOrBurn = 1
// 	} else if toSave.To == chain.EmptyAddressInHex {
// 		toSave.TransferOrMintOrBurn = 2
// 	}
// 	return DB.Create(&toSave).Error
// }

// func SaveEventNekoSpiritUpgrade(event *chain.SpiritUpgrade) error {
// 	toSave := EventNekoSpiritUpgrade{
// 		Hash: Hash{
// 			TransactionHash: event.Raw.TxHash.Hex(),
// 			BlockNumber:     event.Raw.BlockNumber,
// 			BlockHash:       event.Raw.BlockHash.Hex()},
// 		From:          event.Sender.Hex(),
// 		TokenId:       event.TokenId.Uint64(),
// 		NewLevel:      event.NewLevel.Uint64(),
// 		NekoCoinCount: decimal.NewFromBigInt(event.NekoCount, -18),
// 		PrismaCount:   decimal.NewFromBigInt(event.PrismCount, -18),
// 	}
// 	return DB.Create(&toSave).Error
// }

// func SaveEventNekoSpiritUpgradeAscend(event *chain.SpiritUpgradeAscend) error {
// 	toSave := EventAscendUpgrade{
// 		Hash: Hash{
// 			TransactionHash: event.Raw.TxHash.Hex(),
// 			BlockNumber:     event.Raw.BlockNumber,
// 			BlockHash:       event.Raw.BlockHash.Hex()},
// 		From:          event.Sender.Hex(),
// 		NewLevel:      event.NewLevel.Uint64(),
// 		NekoCoinCount: decimal.NewFromBigInt(event.NekoCount, -18),
// 		PrismaCount:   decimal.NewFromBigInt(event.Prism, -18),
// 	}
// 	return DB.Create(&toSave).Error
// }

func CreateAddress(address *ServerAddress) (uint64, error) {
	if err := DB.Create(&address).Error; err != nil {
		return 0, err
	}
	return address.ID, nil
}

// func SaveEventShardTransfer(event *chain.ShardTransfer) error {
// 	toSave := EventTemporalShardTransfer{
// 		Hash: Hash{
// 			TransactionHash: event.Raw.TxHash.Hex(),
// 			BlockNumber:     event.Raw.BlockNumber,
// 			BlockHash:       event.Raw.BlockHash.Hex()},
// 		From:    event.From.Hex(),
// 		To:      event.To.Hex(),
// 		TokenId: event.TokenId.Uint64(),
// 	}
// 	if toSave.From == chain.EmptyAddressInHex {
// 		toSave.TransferOrMintOrBurn = 1
// 	} else if toSave.To == chain.EmptyAddressInHex {
// 		toSave.TransferOrMintOrBurn = 2
// 	}
// 	return DB.Create(&toSave).Error
// }

func CreateShardRecord(record ServerTemporalShardRecord, address string) error {
	Cache.Delete(CacheTagUid + strconv.FormatUint(GetAddressDetailByAddress(address).Uid, 10))
	return DB.Create(&record).Error
}

func DeleteShardRecord(tokenId uint64) error {
	return DB.Where("token_id = ?", tokenId).Delete(&ServerTemporalShardRecord{}).Error
}

func GetIndexerHeight() uint64 {
	var record IndexerRecord
	if err := DB.First(&record).Error; err != nil {
		return 0
	} else {
		return record.BlockHeight
	}
}

func UpdateHeight(height uint64) {
	DB.Model(&IndexerRecord{}).Where("id = 1").Update("block_height", height)
}

// func MakeSureThereIsHeight(height uint64) {
// 	DB.FirstOrCreate(&IndexerRecord{}, IndexerRecord{BlockHeight: height})
// }

func QueryUidThatStakeGreatNeko() []uint64 {
	var uid []uint64
	DB.Model(&ServerNekoSpiritInfo{}).Where("level = ?", 13).Where("is_staked = ?", true).Distinct("stake_from_uid").Select("stake_from_uid").Find(&uid)
	return uid
}

func QueryNotOpenedChest(uid uint64) ServerChest {
	var chest ServerChest
	todayStart := time.Now().Truncate(24 * time.Hour)
	todayEnd := todayStart.Add(24 * time.Hour)
	DB.Where("uid = ?", uid).Where("is_open = ?", false).Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).First(&chest)
	return chest
}

func UpdateChest(chest *ServerChest) {
	DB.Model(&chest).Updates(chest)
}

func CanEmpowerChest(uid uint64) bool {
	var record ServerChestEmpowerRecord
	todayStart := time.Now().Truncate(24 * time.Hour)
	todayEnd := todayStart.Add(24 * time.Hour)
	DB.Where("uid = ?", uid).Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).First(&record)
	return record.ID == 0
}

func CreateEmpowerChestRecord(record ServerChestEmpowerRecord) {
	DB.Create(&record)
}

func CountEmpower(cid uint64) uint64 {
	var count int64
	todayStart := time.Now().Truncate(24 * time.Hour)
	todayEnd := todayStart.Add(24 * time.Hour)
	DB.Model(&ServerChestEmpowerRecord{}).Where("cid = ?", cid).Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).Count(&count)
	return uint64(count)
}

func QueryInviterOfInviter(uid uint64) uint64 {
	var inviter ServerAddress
	DB.Where("id = ?", uid).First(&inviter)
	return inviter.ID
}

func CreateInvitationRecord(record ServerInvitationRecord) error {
	return DB.Create(&record).Error
}

func QueryStakedSpiritList() []ServerNekoSpiritInfo {
	var list []ServerNekoSpiritInfo
	DB.Where("is_staked = ?", true).Find(&list)
	return list
}

func GetBountyWaveList() ([]ServerWhiteListOfBountyWave, decimal.Decimal, bool) {
	now := time.Now()
	var config ServerBountyWaveConfig

	if err := DB.First(&config).Error; err != nil || config.StartTime.Compare(now) > 0 || config.EndTime.Compare(now) < 0 {
		fmt.Println("GetBountyWaveList error: ", err)
		return nil, decimal.Decimal{}, false
	}

	var list []ServerWhiteListOfBountyWave
	DB.Find(&list)
	return list, config.Boost, true
}

func GetRewardPool() decimal.Decimal {
	var result ServerRewardPool
	DB.First(&result)
	return result.TotalAmount
}

func UpdateNekoSpiritList(toUpdate []ServerNekoSpiritInfo) {
	fmt.Println("[Database] UpdateNekoSpiritList size: ", len(toUpdate))
	now := time.Now()
	for _, info := range toUpdate {
		if err := DB.Model(&info).Updates(info).Error; err != nil {
			fmt.Println("UpdateNekoSpiritList error: ", err)
		}
		Cache.Delete(CacheTagNekoSpirit + strconv.FormatUint(info.TokenId, 10))
	}
	fmt.Println("[Database] UpdateNekoSpiritList time: ", time.Since(now).Abs().Milliseconds(), "ms")
}

func CreateInvitationRewardRecords(records []ServerInvitationRewardRecord) error {
	return DB.Create(&records).Error
}

func CreateInvitationRewardStatic(uid uint64) {
	DB.Create(&ServerInvitationRewardStatic{Uid: uid})
}

func AddInvitationRewardStatic(uid uint64, amount decimal.Decimal) {
	DB.Model(&ServerInvitationRewardStatic{}).Where("uid = ?", uid).Update("total_amount", gorm.Expr("total_amount + ?", amount))
}

func UnlockInvitationRewardStatic(uid uint64, amount decimal.Decimal) {
	DB.Model(&ServerInvitationRewardStatic{}).Where("uid = ?", uid).Update("unlocked_amount", gorm.Expr("unlocked_amount + ?", amount))
}

func ClaimInvitationRewardStatic(uid uint64, amount decimal.Decimal) {
	DB.Model(&ServerInvitationRewardStatic{}).Where("uid = ?", uid).Update("claimed_amount", gorm.Expr("claimed_amount + ?", amount))
}

func CreateClaimRecord(uid uint64, amount decimal.Decimal) {
	DB.Create(&ServerClaimNekoSpiritRecord{Uid: uid, Amount: amount})
}

func AddRewardsPool(amount decimal.Decimal) {
	DB.Model(&ServerRewardPool{}).Where("id = 1").Update("total_amount", gorm.Expr("total_amount + ?", amount))
}

func SubRewardPool(amount decimal.Decimal) {
	DB.Model(&ServerRewardPool{}).Where("id = 1").Update("total_amount", gorm.Expr("total_amount - ?", amount))
}

func QueryStarterChestConfig() ServerStarterChestConfig {
	var config ServerStarterChestConfig
	DB.First(&config)
	return config
}

func AddStarterChestOpened() {
	DB.Model(&ServerStarterChestConfig{}).Where("id = 1").Update("opened", gorm.Expr("opened + ?", 1))
}

func UpdateAddressStarter(address string) {
	DB.Model(&ServerAddress{}).Where("address = ?", address).Update("is_starter", false)
}

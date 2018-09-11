package rediscache

const basic = "basic"
const standard = "standard"
const premium = "premium"

type planDetail struct {
	planName          string
	allowedCapacity   []int64
	allowedShardCount []int64
}

func newBasicPlanDetail() planDetail {
	return planDetail{
		planName:        basic,
		allowedCapacity: []int64{0, 1, 2, 3, 4, 5, 6},
	}
}

func newStandardPlanDetail() planDetail {
	return planDetail{
		planName:        standard,
		allowedCapacity: []int64{0, 1, 2, 3, 4, 5, 6},
	}
}

func newPremiumPlanDetail() planDetail {
	return planDetail{
		planName:          premium,
		allowedCapacity:   []int64{1, 2, 3, 4},
		allowedShardCount: []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}
}

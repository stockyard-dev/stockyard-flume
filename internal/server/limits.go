package server

type Tier string

const (
	TierFree Tier = "free"
	TierPro  Tier = "pro"
)

type Limits struct {
	Tier        Tier
	Description string
}

func LimitsFor(tier string) Limits {
	if tier == "pro" {
		return Limits{Tier: TierPro, Description: "Unlimited sources, 10M logs"}
	}
	return Limits{Tier: TierFree, Description: "2 sources, 50k logs retained"}
}

func (l Limits) IsPro() bool {
	return l.Tier == TierPro
}

package goconf

import "time"

type TableService struct {
	Table
	Scheduler
	GKE
	Multi
	Cleaning
	Enrich
	Daily
	Tournament
}

type Table struct {
	SeatReservationTimeout      time.Duration   `mapstructure:"seat_reservation_timeout"`
	MinDecisionTimeout          time.Duration   `mapstructure:"min_decision_timeout"`
	MaxDecisionTimeout          time.Duration   `mapstructure:"max_decision_timeout"`
	GameEndMinTimeout           time.Duration   `mapstructure:"game_end_min_timeout"`
	GameEndPotTimeout           time.Duration   `mapstructure:"game_end_pot_timeout"`
	GameEndCommunityCardTimeout time.Duration   `mapstructure:"game_end_community_card_timeout"`
	ShowDownTimeout             time.Duration   `mapstructure:"show_down_timeout"`
	FinishedLobbies             FinishedLobbies `mapstructure:"finished_lobbies"`
	PlayerActionDuration        time.Duration   `mapstructure:"player_action_duration"`
	RakePercent                 float64         `mapstructure:"rake_percent"`
	MaxRake                     int64           `mapstructure:"max_rake"`
}

type Tournament struct {
	FeePercent float64 `mapstructure:"fee_percent"`
}

type Scheduler struct {
	Cron
}

type Cron struct {
	TimeoutRecovery                   string `mapstructure:"timeout_recovery"`
	StartMultiTournaments             string `mapstructure:"start_multi_tournaments"`
	CleanSittingOutPlayers            string `mapstructure:"clean_sitting_out_players"`
	CleanWaitingTables                string `mapstructure:"clean_waiting_tables"`
	CleanAlonePlayerOnPersistentTable string `mapstructure:"clean_alone_player_on_persistent_table"`
	DeleteFinishedLobbies             string `mapstructure:"delete_finished_lobbies"`
	DeleteNotStartedSitngo            string `mapstructure:"delete_not_started_sitngo"`
	CreateDailySitngo                 string `mapstructure:"create_daily_sitngo"`
}

type Multi struct {
	RebalancerPeriod time.Duration `mapstructure:"rebalancer_period"`
	TableSize        int           `mapstructure:"table_size"`
	DecisionTimeout  time.Duration `mapstructure:"decision_timeout"`
	Freerolls
}

type Daily struct {
	LastVideoID string `mapstructure:"last_video_id"`
}

type Freerolls struct {
	At []string `mapstructure:"at"`
}

type Cleaning struct {
	TournamentSittingOutTimeout  time.Duration `mapstructure:"tournament_sitting_out_players_timeout"`
	CashSittingOutTimeout        time.Duration `mapstructure:"cash_sitting_out_players_timeout"`
	WaitingTablesTimeout         time.Duration `mapstructure:"waiting_tables_timeout"`
	AlonePlayerOnPersistentTable time.Duration `mapstructure:"alone_player_on_persistent_table"`
	SitngoStartTimeout           time.Duration `mapstructure:"sitngo_start_timeout"`
}

type FinishedLobbies struct {
	Period time.Duration
}

type Enrich struct {
	PlayersEnabled bool `mapstructure:"players_enabled"`
}

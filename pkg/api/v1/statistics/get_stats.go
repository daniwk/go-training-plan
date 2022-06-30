package statistics

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/daniwk/training-plan/pkg/models"
	"github.com/gin-gonic/gin"
)

type OverallActivities struct {
	PlannedActivities int `json:"planned_activities"`
	ActualActivities  int `json:"actual_activities"`
	PlannedDuration   int `json:"planned_duration"`
	ActualDuration    int `json:"actual_duration"`
}

type Distance struct {
	PlannedKm                             int     `json:"planned_km"`
	PlannedDistanceIncreaseFromLastWeek   float64 `json:"planned_distance_increase_vs_last_week"`
	PlannedPercentageIncreaseFromLastWeek float64 `json:"planned_percentage_increase_vs_last_week"`
	ActualKm                              float64 `json:"actual_km"`
	ActualDistanceIncreaseFromLastWeek    float64 `json:"actual_distance_increase_vs_last_week"`
	ActualPercentageIncreaseFromLastWeek  float64 `json:"actual_percentage_increase_vs_last_week"`
	AcuteCronicRatio                      float64 `json:"acute_cronic_ratio"`
}

type Workouts struct {
	PlannedNumber             int     `json:"planned_number"`
	ActualNumber              int     `json:"actual_number"`
	PlannedMinutesInQuality   int     `json:"planned_minutes_in_quality"`
	ActualMinutesInQuality    float64 `json:"actual_minutes_in_quality"`
	TotalSufferScore          float64 `json:"total_suffer_score"`
	AverageHeartRateInWorkout float64 `json:"avg_heart_rate_in_workout"`
	MaxHeartRateInWorkout     float64 `json:"max_heart_rate_in_workout"`
	AverageSpeedInWorkoutsKmh float64 `json:"average_speed_in_workout_kmh"`
	MaxSpeedInWorkoutsKmh     float64 `json:"max_speed_in_workout_kmh"`
}

type Intensity struct {
	PlannedLoad       float64 `json:"planned_load"`
	ActualLoad        float64 `json:"actual_load"`
	StravaSufferScore float64 `json:"strava_suffer_score"`
}

type GenericWeeklyStatitistics struct {
	Activities OverallActivities `json:"activities"`
	Distance   Distance          `json:"distance"`
	Workouts   Workouts          `json:"workouts"`
	Intensity  Intensity         `json:"intensity"`
}

type WeeklyStatitistics struct {
	StartDate      time.Time                 `json:"start_date"`
	EndDate        time.Time                 `json:"end_date"`
	Running        GenericWeeklyStatitistics `json:"running"`
	WeightTraining GenericWeeklyStatitistics `json:"weight_Training"`
	Rides          GenericWeeklyStatitistics `json:"rides"`
}

type ActivityLoad struct {
	PlannedLoad float64
	ActualLoad  float64
}

type WorkoutLoad struct {
	MinutesInQuality          float64
	AverageHeartRateInWorkout float64
	MaxHeartRateInWorkout     float64
	AverageSpeedInWorkoutsKmh float64
	MaxSpeedInWorkoutsKmh     float64
}

// Iterate over workout's laps and calculate MinutesInQuality and AverageHeartRate for workout
func CalculateWorkoutLoad(laps []models.Lap) WorkoutLoad {
	var actual_seconds_in_quality int
	var counter int
	var heartrate_in_quality float64
	var max_heart_rate_in_workout float64
	var speed_in_quality float64
	var max_speed_in_quality float64
	const minimum_workout_speed = 3.63 // 3.63 m/s ~ 4:35 min/km

	for _, lap := range laps {

		// Only account for laps with greater pace than 4:35 min/km ~ 3.63 meter/second (which is unit Strava API uses in max_speed and average_speed)
		if lap.AverageSpeed > minimum_workout_speed {

			counter += 1
			actual_seconds_in_quality += lap.MovingTime
			heartrate_in_quality += lap.AverageHeartRate
			if max_heart_rate_in_workout < lap.MaxHeartRate {
				max_heart_rate_in_workout = lap.MaxHeartRate
			}
			speed_in_quality += lap.AverageSpeed
			if max_speed_in_quality < lap.MaxSpeed {
				max_speed_in_quality = lap.MaxSpeed
			}
		}
	}

	response := WorkoutLoad{}
	response.MinutesInQuality = float64(actual_seconds_in_quality / 60)
	response.AverageHeartRateInWorkout = heartrate_in_quality / float64(counter)
	response.MaxHeartRateInWorkout = max_heart_rate_in_workout
	response.AverageSpeedInWorkoutsKmh = speed_in_quality / float64(counter) * 3.6
	response.MaxSpeedInWorkoutsKmh = max_speed_in_quality * 3.6

	return response
}

// Calculate planned and actual activity load, i.e. create own SufferScore.
func CalculateActivityLoad(planned_activity models.PlannedActivity) ActivityLoad {

	activity_load := ActivityLoad{}
	var planned_power = 1
	if planned_activity.Intensity < 8 && planned_activity.Intensity > 3 {
		planned_power += 1
	} else if planned_activity.Intensity > 7 {
		planned_power += 2
	}
	activity_load.PlannedLoad = float64(planned_activity.Intensity) * float64(planned_activity.Duration) * float64(planned_power)

	if planned_activity.StravaActivity != nil {
		var actual_power = 1
		if planned_activity.StravaActivity.PerceivedExertion < 8 && planned_activity.StravaActivity.PerceivedExertion > 3 {
			actual_power += 1
		} else if planned_activity.StravaActivity.PerceivedExertion > 7 {
			actual_power += 2
		}
		activity_load.ActualLoad = planned_activity.StravaActivity.PerceivedExertion * float64(planned_activity.StravaActivity.MovingTime) / 60 * float64(actual_power)
	}

	return activity_load
}

type PeriodStatistics struct {
	WeeklyStatitistics WeeklyStatitistics
	Error              error
}

func GetStatisticsForPeriod(h handler, start_date time.Time, end_date time.Time) PeriodStatistics {
	var weekly_stats WeeklyStatitistics
	var period_statistics PeriodStatistics
	weekly_stats.StartDate = start_date
	weekly_stats.EndDate = end_date
	var planned_activities []models.PlannedActivity

	if result := h.DB.Model(&models.PlannedActivity{}).Preload("StravaActivity").Preload("StravaActivity.Laps").Where("activity_type = ? AND date > ? AND date < ?", models.Run, start_date, end_date).Find(&planned_activities); result.Error != nil {
		period_statistics.Error = errors.New("NotFound")
		return period_statistics
	}

	// Read all activities for the gived week

	for _, planned_activity := range planned_activities {
		fmt.Printf("Planned activity %v\n", planned_activity)
		weekly_stats.Running.Activities.PlannedActivities += 1
		weekly_stats.Running.Activities.PlannedDuration += planned_activity.Duration
		weekly_stats.Running.Distance.PlannedKm += planned_activity.Distance

		activity_load := CalculateActivityLoad(planned_activity)
		weekly_stats.Running.Intensity.PlannedLoad += activity_load.PlannedLoad

		if planned_activity.WorkoutType == models.WORKOUT {
			weekly_stats.Running.Workouts.PlannedNumber += 1
			weekly_stats.Running.Workouts.PlannedMinutesInQuality += planned_activity.MinutesInQuality
		}

		if planned_activity.StravaActivity != nil {
			weekly_stats.Running.Activities.ActualActivities += 1
			weekly_stats.Running.Activities.ActualDuration += planned_activity.StravaActivity.MovingTime / 60
			weekly_stats.Running.Distance.ActualKm += float64(planned_activity.StravaActivity.Distance) / 1000
			weekly_stats.Running.Intensity.ActualLoad += planned_activity.StravaActivity.PerceivedExertion
			weekly_stats.Running.Intensity.StravaSufferScore += planned_activity.StravaActivity.SufferScore
			weekly_stats.Running.Intensity.ActualLoad += activity_load.ActualLoad

			if planned_activity.WorkoutType == models.WORKOUT {

				workout_load := CalculateWorkoutLoad(planned_activity.StravaActivity.Laps)

				weekly_stats.Running.Workouts.ActualNumber += 1
				weekly_stats.Running.Workouts.TotalSufferScore += planned_activity.StravaActivity.SufferScore
				weekly_stats.Running.Workouts.ActualMinutesInQuality += workout_load.MinutesInQuality
				weekly_stats.Running.Workouts.AverageHeartRateInWorkout += workout_load.AverageHeartRateInWorkout
				weekly_stats.Running.Workouts.MaxHeartRateInWorkout = workout_load.MaxHeartRateInWorkout
				weekly_stats.Running.Workouts.AverageSpeedInWorkoutsKmh = workout_load.AverageSpeedInWorkoutsKmh
				weekly_stats.Running.Workouts.MaxSpeedInWorkoutsKmh = workout_load.MaxSpeedInWorkoutsKmh
			}
		}
	}
	period_statistics.WeeklyStatitistics = weekly_stats
	return period_statistics
}

func (h handler) GetWeeklyRunStatistics(c *gin.Context) {
	start_day := c.Query("start_day")
	start_day_int, err := strconv.Atoi(start_day)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	start_month := c.Query("start_month")
	start_month_int, err := strconv.Atoi(start_month)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	start_year := c.Query("start_year")
	start_year_int, err := strconv.Atoi(start_year)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	start_date := time.Date(int(start_year_int), time.Month(start_month_int), start_day_int, 0, 0, 0, 0, time.Local)
	end_date := start_date.Add(time.Hour * 24 * 7)

	// Read all activities for the gived week
	stats := GetStatisticsForPeriod(h, start_date, end_date)
	if stats.Error != nil {
		if errors.Is(stats.Error, errors.New("NotFound")) {
			c.AbortWithError(http.StatusNotFound, stats.Error)
		} else {
			c.AbortWithError(http.StatusInternalServerError, stats.Error)
		}
		return
	}

	// Calulcate progress vs last week and vs average last 4 weeks
	start_date_minus_1_week := start_date.Add(time.Hour * 24 * -7)
	last_week_stats := GetStatisticsForPeriod(h, start_date_minus_1_week, start_date)

	fmt.Printf("Last week's stats: %v\n", last_week_stats.WeeklyStatitistics)
	stats.WeeklyStatitistics.Running.Distance.PlannedDistanceIncreaseFromLastWeek = float64(stats.WeeklyStatitistics.Running.Distance.PlannedKm) - last_week_stats.WeeklyStatitistics.Running.Distance.ActualKm
	// stats.WeeklyStatitistics.Running.Distance.PlannedPercentageIncreaseFromLastWeek = (float64(stats.WeeklyStatitistics.Running.Distance.PlannedKm) - last_week_stats.WeeklyStatitistics.Running.Distance.ActualKm) / last_week_stats.WeeklyStatitistics.Running.Distance.ActualKm

	fmt.Printf("stats: %v\n", stats.WeeklyStatitistics)
	c.JSON(http.StatusOK, &stats.WeeklyStatitistics)
}

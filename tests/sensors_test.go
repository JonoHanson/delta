package meta_test

import (
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestSensors(t *testing.T) {
	var sensors meta.InstalledSensors

	t.Log("Load installed sensors file")
	{
		if err := meta.LoadLists("../installs", "sensors.csv", &sensors); err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check for sensor installation location overlaps")
	{
		installs := make(map[string]meta.InstalledSensors)
		for _, s := range sensors {
			_, ok := installs[s.Station]
			if ok {
				installs[s.Station] = append(installs[s.Station], s)

			} else {
				installs[s.Station] = meta.InstalledSensors{s}
			}
		}

		var keys []string
		for k, _ := range installs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := installs[k]

			for i, n := 0, len(v); i < n; i++ {
				for j := i + 1; j < n; j++ {
					switch {
					case v[i].Location != v[j].Location:
					case v[i].EndTime.Before(v[j].StartTime):
					case v[i].StartTime.After(v[j].EndTime):
					case v[i].EndTime.Equal(v[j].StartTime):
					case v[i].StartTime.Equal(v[j].EndTime):
					default:
						t.Errorf("sensor %s/%s at %-5s has location %-2s overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Station, v[i].Location, v[i].StartTime.Format(meta.DateTimeFormat), v[i].EndTime.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	}

	t.Log("Check for sensor installation equipment overlaps")
	{
		installs := make(map[string]meta.InstalledSensors)
		for _, s := range sensors {
			_, ok := installs[s.Model]
			if ok {
				installs[s.Model] = append(installs[s.Model], s)

			} else {
				installs[s.Model] = meta.InstalledSensors{s}
			}
		}

		var keys []string
		for k, _ := range installs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := installs[k]

			for i, n := 0, len(v); i < n; i++ {
				for j := i + 1; j < n; j++ {
					switch {
					case v[i].Serial != v[j].Serial:
					case v[i].EndTime.Before(v[j].StartTime):
					case v[i].StartTime.After(v[j].EndTime):
					case v[i].EndTime.Equal(v[j].StartTime):
					case v[i].StartTime.Equal(v[j].EndTime):
					default:
						t.Errorf("sensor %s/%s at %-5s has location %-2s overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Station, v[i].Location, v[i].StartTime.Format(meta.DateTimeFormat), v[i].EndTime.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	}

	t.Log("Check for missing sensor stations")
	{
		var stations meta.Stations

		if err := meta.LoadList("../network/stations.csv", &stations); err != nil {
			t.Fatal(err)
		}

		keys := make(map[string]interface{})

		for _, s := range stations {
			keys[s.Code] = true
		}

		for _, s := range sensors {
			if _, ok := keys[s.Station]; ok {
				continue
			}
			t.Errorf("unable to find sensor installed station %-5s", s.Station)
		}
	}

}

package timer 
import ("time")


func get_wall_time() float64 {
    var time struct {
        tv_sec  int64
        tv_usec int64  
    }
    gettimeofday(&time, nil)
    return float64(time.tv_sec) + float64(time.tv_usec)*0.000001
}

float64 timerEndTime;
bool timerActive;

func timer_start(duration float64) {
    timerEndTime = get_wall_time() + duration
    timerActive = true
}

func timer_stop() {
    timerActive = false
}

func timer_timedOut() bool {
    return timerActive && get_wall_time() > timerEndTime 
}

static  double          timerEndTime;
static  int             timerActive;

void timer_start(double duration){
    timerEndTime    = get_wall_time() + duration;
    timerActive     = 1;
}

void timer_stop(void){
    timerActive = 0;
}

int timer_timedOut(void){
    return (timerActive  &&  get_wall_time() > timerEndTime);
}




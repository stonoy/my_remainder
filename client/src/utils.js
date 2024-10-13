import axios from "axios";

export const axiosBase = axios.create({
    baseURL : "http://localhost:8080/api/v1"
})

export const extractDayTime = (timimg) => {
    const dayTime = {
        day: 0,
        hour: 0,
        miniute: 0,
        second: 0
    }

    const now = new Date()
    const time = new Date(timimg.toLocaleString())

    if (now > time) {
        return dayTime
    }

    const diff = time - now

    const day = 1000*60*60*24
    const hour = 1000*60*60
    const miniute = 1000*60
    const second = 1000

    dayTime.day = Math.floor(diff / day)
    dayTime.hour = Math.floor((diff%day)/hour)
    dayTime.miniute = Math.floor((diff%hour)/miniute)
    dayTime.second = Math.floor((diff%miniute)/second)

    return dayTime
}

export const ticking = (remainder) => {
    let {day, hour, miniute, second} = remainder

    if (second > 0){
        second -= 1
    } else if (miniute > 0){
        miniute -= 1
        second = 59
    } else if (hour > 0){
        hour -= 1
        miniute = 59
        second = 59
    } else if (day > 0) {
        day -= 1
        hour = 23
        miniute = 59
        second = 59
    }

    return {...remainder, day, hour, miniute, second}
}
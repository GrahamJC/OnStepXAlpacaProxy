import { ref } from 'vue'

export const onstepx: OnstepX = ref({});

type OnstepX = {
    firmwareVersion: string,
    localDateTime: Date,
    utcOffset: number,
    siteLongitude: number,
    siteLatitude: number,
    siteElevation: number,
    isAtHome: boolean,
    isParked: boolean,
    isSlewing: boolean,
    isGuiding: boolean,
    isTracking: boolean,
    autoHomeAtStartup: boolean,
    raHomeOffset: number,               // 
    decHomeOffset: number,              // 
    currentRA: number,                  // 0-24 hr
    currentDec: number,                 // -90 to +90  deg
    currentAlt: number,                 // -90 to +90 deg
    currentAz: number,                  // 0-360 deg
    targetRA: number,                   // 0-24 hr
    targetDec: number,                  // -90 to +90 deg
    maxSlewSpeed: number,               // Degrees per second
    gotoRate: number,                   // VerySlow, Slow, Normal, Fast or VeryFast
    gotoAlert: boolean,                 // Audible alert when goto complete
    flipAuto: boolean,                  // Auto meridian flip when limit reached
    flipPauseHome: boolean,             // Pause at Home during meridian flip
    flipPierSide: string,               // East, West or Both
    pulseGuideRate: number,             // 0.25, 0.5 or 1 (x sidereal rate)
    guideRate: string,                  // 2x, 4x, 8x, 20x, 48x, Fast or VeryFast
    trackingRate: string,               // Sidereal, Lunar, Solar or King
    trackingCompensation: string,       // Full, RefractionOnly or Off
    trackingDualAxis: boolean,
    backlashRa: number,                 // 0-3600 arcsec
    backlashDec: number,                // 0-3600 arcsec
    limitHorizon: number,               // -30 to +30 deg
    limitOverhead: number,              // 60-90 deg
    limitMeridianEast: number,          // -270 to +270 deg
    limitMeridianWest: number,          // -270 to +270 deg
}
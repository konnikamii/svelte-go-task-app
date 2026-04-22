export type TourStep = {
  element?: string
  message: string
  title?: string
  placement?: 'auto' | 'top' | 'bottom'
}

type TourState = {
  index: number
  open: boolean
  steps: TourStep[]
}

const tour = $state<TourState>({
  index: 0,
  open: false,
  steps: [],
})

export function startTour(steps: TourStep[]) {
  if (steps.length === 0) {
    stopTour()
    return
  }

  tour.steps = steps.map((step) => ({
    placement: 'auto',
    ...step,
  }))
  tour.index = 0
  tour.open = true
}

export function stopTour() {
  tour.open = false
  tour.index = 0
  tour.steps = []
}

export function nextTourStep() {
  if (tour.index >= tour.steps.length - 1) {
    stopTour()
    return
  }

  tour.index += 1
}

export function prevTourStep() {
  if (tour.index === 0) return
  tour.index -= 1
}

export function useTour() {
  return {
    nextTourStep,
    prevTourStep,
    stopTour,
    tour,
  }
}

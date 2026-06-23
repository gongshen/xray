import {
  bindEmitterHandlers as bindSharedEmitterHandlers,
  bindWindowEvent,
} from '../../../../utils/eventLifecycle.mjs'

export function bindBodyClickHandler(handler, target = document.body) {
  return bindWindowEvent(target, 'click', handler)
}

export function bindEmitterHandlers(emitter, handlers) {
  return bindSharedEmitterHandlers(emitter, handlers)
}

import { EntityActionTypes, actionCreators, EntityAction, SelectEntity, MovingEntity } from './actions'
import { actionCreators as controlActionCreators } from 'modules/control/actions'
import { actionCreators as relationActionCreators } from 'modules/relation/actions'
import { Reducer } from 'redux'
import { call, put, takeEvery, takeLatest, take, select } from 'redux-saga/effects'
import { newRelation } from 'grpc/util/relation'

import { RelAssociation } from 'grpc/erd_pb'

import { List } from 'immutable'
import { RootState } from 'modules/rootReducer'

function* onCreateEntity(action: EntityAction) {
    const state = <RootState>(yield select())
    if (!state.control.repeatAction) {
        yield put(controlActionCreators.finishCreateEntity())
    }
}


function* onSelectEntity(action: EntityAction) {
    const currentState = <RootState>(yield select())
    const selectAction = <SelectEntity>action

    const objectId = selectAction.payload

    // if add relation
    if (currentState.control.connectingOneToMenyRel) {
        if (currentState.entity.seqentialChoiceEntities.contains(objectId)) {
            return
        }

        yield put(actionCreators.seqChoiceEntities(objectId))

        const latestState = <RootState>(yield select())

        console.debug("choice")
        console.debug(latestState.entity.seqentialChoiceEntities.size)

        if (latestState.entity.seqentialChoiceEntities.size >= 2) {
            //TODO (taco) add rel
            const beginEntityObjId = latestState.entity.seqentialChoiceEntities.get(0)
            const beginEntity = latestState.entity.entities.get(beginEntityObjId)

            const endEntityObjId = latestState.entity.seqentialChoiceEntities.get(1)
            const endEntity = latestState.entity.entities.get(endEntityObjId)

            const newRel = newRelation(
                beginEntity,
                endEntity,
                RelAssociation.One,
                RelAssociation.Many)

            yield put(relationActionCreators.addRelation(newRel))
            yield put(actionCreators.cancelSelection())
            yield put(controlActionCreators.finishConnectingRel())
        }

    } else {
        yield put(actionCreators.choiceEntities(objectId))
    }
}

function* onMovingEntity(action: EntityAction) {
    const currentState = <RootState>(yield select())
    const movingEntity = <MovingEntity>action
    const { objectId, coord } = movingEntity.payload
    yield put(relationActionCreators.rerenderByEntityMove(objectId, coord))
}

export function* entityWatcher() {
    yield takeLatest(EntityActionTypes.CREATE_NEW_ENTITY, onCreateEntity)
    yield takeLatest(EntityActionTypes.SELECT_ENTITY, onSelectEntity)
    yield takeLatest(EntityActionTypes.MOVING_ENTITY, onMovingEntity)
}


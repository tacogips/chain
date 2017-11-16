import * as React from 'react'
import { render } from 'react-dom'
import { Provider } from 'react-redux'
import { createStore, applyMiddleware } from 'redux'
import createSagaMiddleware from 'redux-saga'
import * as Konva from 'konva'
import { Stage, Layer, Rect, Group, Circle } from 'react-konva'
import { Map } from 'immutable'

import { Entity, Rel ,Coord} from 'grpc/erd_pb'

export interface EntityPanelProps {
    key: string
    entity: Entity
}

export interface EntityPanelProps {
	dragStartPos: Coord
}

export class EntityPanel extends React.Component<EntityPanelProps, {}>{
    //TODO(tacogips) if defined this type as "Reat", Konva.Rect methods will not found? check out after react 16.1 and its typedefinision released
    private refRect: any

    constructor(props?: EntityPanelProps, context?: any) {
        super(props, context)
    }

    changeSize() {

    }

    onMouseOver = () => {
        document.body.style.cursor = 'pointer';
    }

    onMouseOut = () => {
        document.body.style.cursor = 'default';
    }

		onDragStart = (evt.any)  =>{
		}

    onDragEnd = (evt:any) => {
			console.log(evt)

    }

    handleClick = () => {
        this.refRect.fill("red")
        this.refRect.draw()
    }

    render() {
        const { entity} = this.props
        return (
            <Rect
                x={entity.getCoord().getX()}
                y={entity.getCoord().getY()}
                width={entity.getWidthHeight().getW()}
                height={entity.getWidthHeight().getH()}
                fill={entity.getColor()}
                shadowBlur={1}
                draggable={true}
                onDragEnd={this.onDragEnd}
                ref={(ref) => this.refRect = ref}
                onMouseOver={this.onMouseOver}
                onMouseOut={this.onMouseOut}
                onClick={this.handleClick}
            />
        );
    }
}



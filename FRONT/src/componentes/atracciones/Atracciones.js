import React from 'react'
import { DataTable } from 'primereact/datatable'
import { Column } from 'primereact/column';

class Atracciones extends React.Component {
    /**
     * Constructor de la clase Atracciones
     * @param {} props 
     */
    constructor(props) {
        super(props)
        this.state = {
            atracciones: [],
            isFetch: true,
        }
    }
    /**
     * ComponentDidMount que carga la información de las atracciones
     */
    componentDidMount() {
        fetch("http://192.168.43.50:8082/atracciones")
            .then(response => response.json())
            .then(visitantesJson => this.setState( {
                visitantes: visitantesJson.data,
                isFetch: false
            }))
            .catch(error => console.log(error))
    }
    /**
     * Render que muestra la información de las atracciones
     * @returns : Renderizado de las atracciones
     */
    render () {
        
        const { visitantes, isFetch } = this.state

        if (isFetch) {
            return <div>El estado de las atracciones no está disponible por el momento</div>
        }

        return (
          <div className ="container">
            <DataTable value={visitantes}>
                <Column field="id" header="ID"></Column>
                <Column field="tciclo" header="Tiempo de ciclo"></Column>
                <Column field="nvisitantes" header="Capacidad de visitantes"></Column>
                <Column field="posicionx" header="Posición x"></Column>
                <Column field="posiciony" header="Posición y"></Column>
                <Column field="tiempoEspera" header="Tiempo de espera en seg"></Column>
                <Column field="estado" header="Estado"></Column>
            </DataTable>
          </div>
        );
    }
    
}
export default Atracciones;
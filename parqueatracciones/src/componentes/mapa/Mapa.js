import React, {useState, useEffect} from 'react'
import axios from 'axios'

class Mapa extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            clima: [],
            isFetch: true,
        }
    }

    componentDidMount() {
        fetch("https://api.openweathermap.org/data/2.5/weather?q=Paris&appid=c3d8572d0046f36f0c586caa0e2e1d23&lang=es&units=metric")
            .then(response => response.json())
            .then(climaJson => this.setState( {
                clima: climaJson.weather,
                isFetch: false
            }))
    }

    render () {
        const { isFetch, clima } = this.state
        if (isFetch) {
            return 'Cargando...'
        }
        return (
          <div className ="container">
            
          </div>
            /*
            <div className="container">
                <li>
                    clima[0].description
                </li>
            </div>
            */
        )
    }
}

export default Mapa;

/*
function useDatos() {
        const [weather, setWeather] = useState([])
        useEffect(() => {
            
        return weather

        this.state.clima.map((climas) => (
                    <li> {climas} </li>)
    }
*/
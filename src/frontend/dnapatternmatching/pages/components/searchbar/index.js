import React, { useState } from 'react';
import styles from './SearchBar.module.css';
import Subheading from '../subheading';
import PredictionCard from './components/predictioncard';

const SearchBar = () => {
    const [input, setInput] = useState("");
    const [filled, setFilled] = useState(false);
    const [result, setResult] = useState([]);

    const handleChange = async (event) => {
        event.preventDefault();
        setInput(event.target.value);
        console.log(input);
    };

    const clearInput = () => {
        setInput("");
    };

    const fetchData = async (event) => {
        event.preventDefault();
        fetch('http://localhost:8080/api/searchdisease/' + input, {
            method: 'POST'
        })
        .then(response => response.json())
        .then(json => {
            console.log(json);
            setResult(JSON.stringify(json));
            setFilled(true);
        })
        .catch(function(res){ console.log(res) });
    }

    const generateCt = (ct) => {
        if (ct) {
            let res = [];
            for (let i = 0; i < ct.length; i++) {
                res.push(
                    <div>
                        <PredictionCard key={i} object={ct[i]} />
                        <br/>
                    </div>
                )
            }
            return res
        }
        else {
            return <div className={styles.notFound}>There's no prediction as queried</div>
        }
    }

    return (
        <div className={styles.search}>
            <Subheading 
                Text="Search Prediction Result"
                Color="black"
            />
            <div className={styles.inputs}>
                <div className={styles.button}>
                    {
                        (input !== '')
                        ?   <img src='/searchbar/searchicon.png' className={styles.searchIcon} onClick={fetchData}/>
                        :   <img src='/searchbar/closeicon.png' className={styles.searchIcon} onClick={clearInput}/>
                    }
                </div>
                <input 
                    type="text" 
                    placeholder= { "Search" }
                    onChange={ handleChange }
                    value={ input } />
            </div>
            {
                (filled)
                ?   <div className={styles.resOuter}>
                        { generateCt(JSON.parse(result)) }
                    </div>
                :   <div className={styles.resOuter}>
                    </div>
            }
        </div>
    )
}

export default SearchBar;
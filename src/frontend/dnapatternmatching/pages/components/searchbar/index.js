import React, { useState } from 'react';
import styles from './SearchBar.module.css';
import Subheading from '../subheading';

const SearchBar = ({Result}) => {
    const [searchedWord, setSearchedWord] = useState("");

    const handleChange = async (event) => {
        const wordInput = event.target.value;
        setSearchedWord(wordInput);
        console.log(searchedWord);
        const response = await fetch('/api/searchdisease', {
            method: 'POST',
            body: JSON.stringify({
                searchedWord,
            }),
            headers: {
                'Content-Type': 'application/json',
            },
        })
        const data = await response.json();
        console.log(data);
    };

    const clearInput = () => {
        setSearchedWord("");
    };

    const generateCt = (ct) => {
        if (ct) {
            let res = []
            for (let i = 0; i < ct.length; i++) {
                res.push(
                    <div key={i}>
                        <span className={styles.resCt}>{ct[i]}</span>
                    </div>
                )
            }
            return res
        }
    }

    return (
        <div className={styles.searchResultContainer}>
            <Subheading 
                Text="Search Prediction Result"
                Color="black"
            />
            <div className={styles.search}>
                <div className={styles.inputs}>
                    <div className={styles.button}>
                        <img src='/searchbar/searchicon.png' className={styles.searchIcon}/>
                        {/* <img src='/searchbar/closeicon.png' className={styles.searchIcon} onClick={clearInput}/> */}
                    </div>
                    <input 
                        type="text" 
                        placeholder= { "Search" }
                        onChange={ handleChange }
                        value={ searchedWord } />
                </div>
                <div className={styles.resOuter}>
                    { generateCt(Result) }
                </div>
            </div>
        </div>
    )
}

export default SearchBar;
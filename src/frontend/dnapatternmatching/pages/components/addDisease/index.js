import React from 'react';
import styles from './AddDisease.module.css';
import Subheading from '../subheading';

const AddDisease = () => {
    return (
        <div className={styles.addDiseaseContainer}>
            <Subheading 
                Text="Add Disease Data"
                Color="black"
            />
            <div className={styles.formDiseaseContainer}>
                <form action="/api/new" method="post" className={styles.formCt}>
                    <label for="roll" className={styles.label} >Disease Name: </label>
                    <input type="text" required className={styles.inputText} />
                    <label for="name" className={styles.label} >DNA Sequence: </label>
                    <input name="logo" type="file" className={styles.inputFile} />
                    <button type="submit" className={styles.submitButton} >Submit</button>
                </form>
            </div>
        </div>
    )
};

export default AddDisease;
import React from 'react';
import styles from './TestDNA.module.css';
import Subheading from '../subheading';

const TestDNA = ({Tanggal, Pengguna, Penyakit, Hasil}) => {
    return (
        <div className={styles.testDNAContainer}>
            <Subheading 
                Text="Test DNA for Disease"
                Color="white"
            />
            <div className={styles.formTestContainer}>
                <form action="/api/new" method="post" className={styles.formCt}>
                    <label className={styles.label} >Username: </label>
                    <input type="text" required className={styles.inputText} />
                    <label className={styles.label} >DNA Sequence: </label>
                    <input name="logo" type="file" className={styles.inputFile} />
                    <label className={styles.label} >Disease Name: </label>
                    <input type="text" required className={styles.inputText} />
                    <button type="submit" className={styles.submitButton} >Submit</button>
                </form>
            </div>
            <div className={styles.resultContainer}>
                <Subheading 
                    Text="Test Result"
                    Color="white"
                />
                <div className={styles.formTestContainer}>
                    <span className={styles.resultText}>{Tanggal}</span>
                    <span className={styles.resultText}> - </span>
                    <span className={styles.resultText}>{Pengguna}</span>
                    <span className={styles.resultText}> - </span>
                    <span className={styles.resultText}>{Penyakit}</span>
                    <span className={styles.resultText}> - </span>
                    <span className={styles.resultText}>{Hasil}</span>
                </div>
            </div>
        </div>  
    )
};

export default TestDNA;
import React, { useState, useEffect } from 'react';
import styles from './TestDNA.module.css';
import Subheading from '../subheading';

const TestDNA = () => {
    const [tanggalprediksi, setTanggalPrediksi] = useState('');
    const [namapasien, setNamaPasien] = useState('');
    const [namapenyakit, setNamaPenyakit] = useState('');
    const [statuspenyakit, setStatusPenyakit] = useState(false);
    const [rantaidna, setRantaiDNA] = useState('');
    const [filled, setFilled] = useState(false);
    const [file, setFile] = useState(null);

    const handleChangeNamaPasien = (e) => {
        setNamaPasien(e.target.value);
        console.log(namapasien);
    }

    const handleChangeNamaPenyakit = (e) => {
        setNamaPenyakit(e.target.value);
        console.log(namapenyakit);
    }

    const handleChangeRantaiDNA = (e) => {
        e.preventDefault();
        setFile(file);
        
        const reader = new FileReader();
        reader.onload = (e) => {
            const text = e.target.result;
            console.log(text);
            setRantaiDNA(text);
        }
        if (e.target.files[0]) {
            reader.readAsText(e.target.files[0]);
        }
        console.log(rantaidna)
    }

    const onClickKMP = async (e) => {
        e.preventDefault();

        fetch("http://localhost:8080/api/getpredictionKMP/" + namapasien + '/' + rantaidna + '/' + namapenyakit,
        {
            method: "POST"
        })
        .then(response => response.json())
        .then(json => {
            console.log(json)
            setTanggalPrediksi(json.tanggalprediksi);
            setStatusPenyakit(json.statuspenyakit); 
            console.log(tanggalprediksi);
            console.log(statuspenyakit);
            setFilled(true);
        })
        .catch(function(res){ console.log(res) });
    }

    const onClickBM = async (e) => {
        e.preventDefault();

        fetch("http://localhost:8080/api/getpredictionBM/" + namapasien + '/' + rantaidna + '/' + namapenyakit,
        {
            method: "POST"
        })
        .then(response => response.json())
        .then(json => {
            console.log(json)
            setTanggalPrediksi(json.tanggalprediksi);
            setStatusPenyakit(json.statuspenyakit); 
            console.log(tanggalprediksi);
            console.log(statuspenyakit);
            setFilled(true);
        })
        .catch(function(res){ console.log(res) });
    }
    
    return (
        <div className={styles.testDNAContainer}>
            <Subheading 
                Text="Test DNA for Disease"
                Color="white"
            />
            <div className={styles.formTestContainer}>
                <form method="post" className={styles.formCt}>
                    <label className={styles.label} >Username: </label>
                    <input type="text" required onChange={handleChangeNamaPasien} className={styles.inputText} />
                    <label className={styles.label} >DNA Sequence: </label>
                    <input name="logo" type="file" onChange={handleChangeRantaiDNA} className={styles.inputFile} />
                    <label className={styles.label} >Disease Name: </label>
                    <input type="text" required onChange={handleChangeNamaPenyakit} className={styles.inputText} />
                    <div className={styles.submitButtonCt}>
                        <button type="submit" onClick={onClickKMP} className={styles.submitButton} >Knuth-Morris-Pratt</button>
                        <div className={styles.spacing}></div>
                        <button type="submit" onClick={onClickBM} className={styles.submitButton} >Boyer-Moore</button>
                    </div>
                </form>
            </div>
            <div className={styles.resultContainer}>
                <Subheading 
                    Text="Test Result"
                    Color="white"
                />
                <div className={styles.formTestContainer}>
                    {
                        (filled && tanggalprediksi !== '')
                        ?   <div>
                                <span className={styles.resultText}>{ tanggalprediksi }</span>
                                <span className={styles.resultText}> — </span>
                                <span className={styles.resultText}>{ namapasien }</span>
                                <span className={styles.resultText}> — </span>
                                <span className={styles.resultText}>{ namapenyakit }</span>
                                <span className={styles.resultText}> — </span>
                                <span className={styles.resultText}>{ statuspenyakit.toString() }</span>
                            </div>
                        :   (filled && tanggalprediksi === '')
                            ?   <div>
                                    <span className={styles.resultText}>Disease Not Found</span>
                                </div>
                            :
                                <div>
                                    <span className={styles.resultText}></span>
                                </div>
                    }
                </div>
            </div>
        </div>  
    )
};

export default TestDNA;
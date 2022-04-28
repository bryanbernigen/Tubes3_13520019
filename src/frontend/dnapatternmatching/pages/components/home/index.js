import React, { useEffect } from 'react';
import styles from './Home.module.css';
import Typewriter from 'typewriter-effect';

const HomePage = () => {
    useEffect(()=>{
        console.log("masuk use effect home page");

        //INI CONTOH HTTP GET
        // fetch('http://localhost:8080')
        // .then(response => response.json())
        // .then(json => {
        //     // set useState terjadinya disini
        //     console.log(json)
        // })
    
        //INI CONTOH HTTP POST
        fetch("http://localhost:8080/ininamapenyakit",
        {
            method: "POST",
            body: JSON.stringify({"name":"jeff"})
        })
        .then(response => response.json())
        .then(json => {
            // set useState terjadinya disini
            console.log(json)
        })
        // .catch(function(res){ console.log(res) })

        //    let bodyContent = JSON.stringify({
        //      "nama": "jeff"
        //    });
           

        //    fetch("http://localhost:8080/abcd", { 
        //      method: "POST",
        //      body: bodyContent,
        //     })
        //     .then(response => response.json())
        //     .then(json => {
        //         // set useState terjadinya disini
        //         console.log(json)
        //     })

    },[])
    return (
        <div className={styles.homeContainer}>
            <div className={styles.homeBG}>
                <div className={styles.homeTitle}>
                    <Typewriter
                        onInit={(typewriter) => {
                            typewriter.typeString('DNA Pattern Matching')
                                .pauseFor(2500)
                                .start()
                                .deleteAll();
                        }}
                        options={{
                            autoStart: true,
                            loop: true,
                        }}
                    />
                </div>
            </div>
        </div>
    )
}

export default HomePage
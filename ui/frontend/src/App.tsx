// import {useState} from 'react';
// import logo from './assets/images/logo-universal.png';
import TitleBar from './components/tittle_bar';
// import {Greet} from "../wailsjs/go/main/App";
// import { Button } from './components/ui/button';


function App() {
    // const [resultText, setResultText] = useState("Please enter your name below 👇");
    // const [name, setName] = useState('');
    // const updateName = (e: any) => setName(e.target.value);
    // const updateResultText = (result: string) => setResultText(result);

    // function greet() {
    //     Greet(name).then(updateResultText);
    // }

    return (
       <div className='dark @container/main flex flex-1 flex-col gap-2 h-screen h-min-full'>
        <TitleBar/>
       </div>
    )
}

export default App

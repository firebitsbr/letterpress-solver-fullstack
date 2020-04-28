import { Component, OnInit } from '@angular/core';
import { Http, Headers } from '@angular/http';
import { env } from '../../../environments/environment';
import * as matchInfo from 'matchInfo';
import 'rxjs/add/operator/map';

@Component({
  selector: 'app-match',
  templateUrl: './match.component.html',
  styleUrls: ['./match.component.css']
})
export class MatchComponent implements OnInit {

  playersId: string[];

  matches: matchInfo.MatchInfo['matches'];
  opponentNames: matchInfo.Participant['userName'][];
  opponentAvatars: matchInfo.Participant['avatarURL'][];
  lastPlayedWords: string[];
  letterGrids: matchInfo.Match['letters'][];
  tileGrids: matchInfo.Tile[][];

  selectedTile: boolean[][];
  foundWords: string[][];

  choosingWord: string[];

  constructor(private http: Http) {
    this.playersId = env.player.map(p => p.id);
  }

  ngOnInit() {
    this.fetchGames();
    window.onscroll = () => {
      if (window.scrollY == 0) {
        this.fetchGames();
      }
    };
  }

  fetchGames() {
    this.http.get('http://' + window.location.host + '/match')
      .map((resp) => resp.text().length > 1000 ? resp.json() : '')
      .subscribe(
        (data) => {
          if (data) {
            this.processGameData(data);
            // auto find words for matches that already started
            for (let i = 0; i < this.matches.length; i++) {
              this.selectedTile[i] = Array<boolean>(25);
              if (this.matches[i].matchStatus === 4) {
                const letters = this.tileGrids[i].map(t => t.t).join('').toUpperCase();
                this.http.get('http://' + window.location.host + '/letterFrequency?letters=' + letters)
                .map((resp) => resp.json())
                .subscribe((data :number[]) => {
                  const alphabet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
                  const freq = {};
                  data.forEach((d,i) => { if (d > 0) freq[alphabet.charAt(i)] = d; });
                  const maxFreq = Math.max.apply(null, data);
                  const minFreq = Math.min.apply(null, data.filter(d => d>0));
                  console.log('letter freq(A-Z)', freq);
                  console.log('max frequency', maxFreq);
                  console.log('min frequency', minFreq);
                  const tg = this.tileGrids[i];
                  for (let k = 0; k < 25; k++) {
                    tg[k].colorCode = `hsl(103, 90%, ${(Math.log1p(freq[tg[k].t]) / Math.log1p(maxFreq)) ** 1 * 100}%`;
                  } 
                })
                continue; //matchStatus==4: new game, escape
              }
              this.selectAllWhite(i);
              this.findWords(i).then(isFound => {
                if (!isFound) this.selectAllPink(i);
              })
            }
          }
        });
  }

  processGameData(data) {
    if (data.matches) {
      // fetch exsisting game list
      this.matches = data.matches;
    } else if (data.match) {
      // fetch newly created game
      this.matches = [data.match];
    } else {
      console.log('Cannot fetch match data from MITM server...');
      return;
    }

    // only show matches on my turn
    this.matches = this.matches.filter(m => this.playersId.includes(m.participants[m.currentPlayerIndex].userId));

    // sort new turns on top
    // this.matches.sort((b, a) => a.participants[a.currentPlayerIndex==1?0:1].turnDate.valueOf()
    //   - b.participants[b.currentPlayerIndex==1?0:1].turnDate.valueOf());
    console.log(this.matches);

    const opponents = this.matches.map(m => this.playersId.includes(m.participants[0].userId) ? m.participants[1] : m.participants[0]);
    this.opponentNames = opponents.map(p => p.userName);
    this.opponentAvatars = opponents.map(p => p.avatarURL ? p.avatarURL : 'https://thesocietypages.org/socimages/files/2009/05/nopic_192.gif');
    this.lastPlayedWords = this.matches.map(m => m.serverData.usedWords).map(ws => ws.length > 0 ? ws[ws.length - 1] : '');

    this.letterGrids = this.matches.map(m => m['letters']);
    console.log(this.letterGrids);
    this.tileGrids = this.matches.map(m => m['serverData']['tiles'])
      .map(t => t.slice(20, 25).concat(t.slice(15, 20).concat(t.slice(10, 15).concat(t.slice(5, 10).concat(t.slice(0, 5))))));
    this.tileGrids.forEach(tg => tg.forEach(t => t.t = t.t.toUpperCase()));


    for (let i = 0; i < this.tileGrids.length; i++) {
      const tg = this.tileGrids[i];
      // reverse owner, if player moves first
      if (this.playersId.includes(this.matches[i].participants[0].userId)) {
        tg.forEach(t => t.o == 1 ? t.o = 0 : t.o == 0 ? t.o = 1 : t.o)
      }

      // set surronded tiles: if every surrounded tiles have same ownership, set true; otherwise, false
      for (let k = 0; k < 25; k++) {
        tg[k].s = [
          ([4, 9, 14, 19, 24].includes(k)) ? undefined : tg[k + 1], // right
          ([0, 5, 10, 15, 20].includes(k)) ? undefined : tg[k - 1], // left
          tg[k + 5], // down
          tg[k - 5]  // up
        ]
          .filter(t => t)
          .every(t => t.o === tg[k].o); // tg[k].o: ownership of adjacent tiles of t; t.o: ownership of t

        tg[k].color = tg[k].o === 127 ? 'white'
          : tg[k].o === 0 && tg[k].s ? 'red'
            : tg[k].o === 0 && !tg[k].s ? 'pink'
              : tg[k].o === 1 && tg[k].s ? 'blue'
                : tg[k].o === 1 && !tg[k].s ? 'azure' : 'error';
      }
    }

    this.selectedTile = Array<boolean[]>(this.matches.length);
    this.foundWords = Array<string[]>(this.matches.length);
    this.choosingWord = Array<string>(this.matches.length);
  }

  //TODO: Remove the recursive call. The second param decides if the recursive call going on. 
  findWords(i: number): Promise<boolean> {
    const letters = this.tileGrids[i].map(t => t.t).join('').toUpperCase();
    let selected = [];
    for (let k = 0; k < 25; k++) {
      if (this.selectedTile[i][k]) {
        selected.push(letters[k])
      }
    }
    console.log(letters);
    console.log(selected.join(''));
    return new Promise(resolve => {
     this.http.get('http://' + window.location.host + '/words?selected=' + selected.join('') + '&letters=' + letters)
      .map(resp => resp.json())
      .subscribe(data => {
        this.foundWords[i] = data;
        const usedWords = this.matches[i].serverData.usedWords;
        // filter out usedWords
        this.foundWords[i] = this.foundWords[i].filter(w => !usedWords.some(uw => uw.indexOf(w.replace('*', '')) === 0));
        if (this.foundWords[i].length == 0) {
          const tg = this.tileGrids[i];
          // Test whether all selected tiles are blank
          const existSelectedNonBlank = this.selectedTile[i].map((b, k) => tg[k].o!==127 && b).some(b => b);
          if (!existSelectedNonBlank) {
            resolve(false);
            return;
          }
          // Recursively unselect letters
          const order = 'JQXZWKVFYBHGMPUDCLTONRAISE';
          let kToUnselect = -1;
          for (let q = 0; q < order.length; q++) {
            const l = order[q];
            for (let k = 0; k < 25; k++) {
              if(this.selectedTile[i][k] && tg[k].t === l && tg[k].o !== 127){
                kToUnselect = k;
                break;
              }
            }
            if (kToUnselect >=0) break;
          }
          this.selectedTile[i][kToUnselect] = false;
          this.findWords(i);
        }
        //TODO: evalue word
        // basic score (-): covers all pink tiles = 0; miss -1
        // aggro score (+): covers white tile; add +1
        // waste score (+): covers blue or dark red tile; add +1
        // critical staus : hard / soft
        // more sophisticated: consider position, maybe need some machine learing
        // <select> the default word <option>
        let offset = 0;
        do {
          this.choosingWord[i] = this.foundWords[i][offset];
          offset++;
        } while (this.choosingWord[i] && this.choosingWord[i].indexOf('*') > 0 && offset < this.foundWords[i].length);

        resolve(true);
      });
    });
  }

  clearSelected(i: number) {
    this.selectedTile[i].fill(false)
  }

  // select all untouched tiles (white)
  selectAllWhite(i: number) {
    const tg = this.tileGrids[i];
    for (let k = 0; k < 25; k++) {
      this.selectedTile[i][k] = (tg[k].o === 127);
    }
  }

  // select all unsurrounded opponent's tiles (pink)
  selectAllPink(i: number) {
    const tg = this.tileGrids[i];
    for (let k = 0; k < 25; k++) {
      this.selectedTile[i][k] = (tg[k].o == 0 && !tg[k].s);
    } 
    this.findWords(i);
  }

  autoClickTiles(i: number) {
    if (this.choosingWord[i] === undefined) return;

    const tg = this.tileGrids[i];
    const letterMap = {}; // key: letter; value: positions[]

    // construct the letterMap, each letter key will have ordered (pink,white...) tiles number as value
    // 1st, push selected tiles
    for (let k = 0; k < 25; k++) {
      if (this.selectedTile[i][k]) {
        if (!letterMap[tg[k].t]) {
          letterMap[tg[k].t] = []
        }
        letterMap[tg[k].t].push(k)
      }
    }
    // 2nd, push unselected tiles, pink first, then white
    ["pink", "white", "red", "azure", "blue"].forEach(color => {
      for (let k = 0; k < 25; k++) {
        if (tg[k].color === color && !this.selectedTile[i][k]) {
          if (!letterMap[tg[k].t]) {
            letterMap[tg[k].t] = []
          }
          letterMap[tg[k].t].push(k)
        }
      }
    })

    // decide the order of tiles to be clicked
    const clickOrderList = []
    for (let idx = 0; idx < this.choosingWord[i].length; idx++) {
      const letter = this.choosingWord[i][idx].toUpperCase();
      if (/^[A-Z]$/.test(letter)) {
        clickOrderList.push(letterMap[letter].shift())
      }
    }

    this.http.post('http://' + window.location.host + '/word?click=' + this.choosingWord[i], clickOrderList)
      .subscribe()
    console.log('autoclicking tiles...', this.choosingWord[i]);
  }

  deleteWord(i: number) {
    this.http.delete('http://' + window.location.host + '/word?delete=' + this.choosingWord[i])
      .subscribe()
    console.log('deleting...', this.choosingWord[i]);
  }

  speakWordUS(word: string) {
    word = word.replace(/\W/gi, '');
    let s = new SpeechSynthesisUtterance(word);
    speechSynthesis.speak(s);
  }
  speakWordUK(word: string) {
    word = word.replace(/\W/gi, '');
    let s = new SpeechSynthesisUtterance(word);
    s.voice = speechSynthesis.getVoices().filter(v => v.lang.indexOf('en-GB') >= 0)[1]
    speechSynthesis.speak(s);
  }
  speakWordFR(word: string) {
    word = word.replace(/\W/gi, '');
    let s = new SpeechSynthesisUtterance(word);
    s.voice = speechSynthesis.getVoices().filter(v => v.lang.indexOf('fr-FR') >= 0)[1]
    speechSynthesis.speak(s);
  }
  speakWordNL(word: string) {
    word = word.replace(/\W/gi, '');
    let s = new SpeechSynthesisUtterance(word);
    s.voice = speechSynthesis.getVoices().filter(v => v.lang.indexOf('nl-NL') >= 0)[0]
    speechSynthesis.speak(s);
  }
}

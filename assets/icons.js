import {dom, library} from '@fortawesome/fontawesome-svg-core';

import {faSpotify} from '@fortawesome/free-brands-svg-icons/faSpotify';
import {faArrowLeft} from '@fortawesome/free-solid-svg-icons/faArrowLeft';
import {faSearch} from '@fortawesome/free-solid-svg-icons/faSearch';
import {faMusic} from '@fortawesome/free-solid-svg-icons/faMusic';
import {faSignOutAlt} from '@fortawesome/free-solid-svg-icons/faSignOutAlt';
import {faSignInAlt} from '@fortawesome/free-solid-svg-icons/faSignInAlt';
import {faQuoteRight} from '@fortawesome/free-solid-svg-icons/faQuoteRight';
import {faPlay} from '@fortawesome/free-solid-svg-icons/faPlay';
import {faPause} from '@fortawesome/free-solid-svg-icons/faPause';
import {faHome} from '@fortawesome/free-solid-svg-icons/faHome';
import {faQuestion} from '@fortawesome/free-solid-svg-icons/faQuestion';
import {faEdit} from '@fortawesome/free-solid-svg-icons/faEdit';

const icons = [
  faArrowLeft, faSearch, faMusic, faSignOutAlt, faSignInAlt,
  faQuoteRight, faPlay, faPause, faHome, faQuestion, faEdit,
  faSpotify,
];

library.add(...icons);
dom.watch();
